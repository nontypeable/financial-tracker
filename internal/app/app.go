package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	customMiddleware "github.com/nontypeable/financial-tracker/internal/app/middleware"
	"github.com/nontypeable/financial-tracker/internal/auth"
	"github.com/nontypeable/financial-tracker/internal/config"
	accountDelivery "github.com/nontypeable/financial-tracker/internal/delivery/account"
	userDelivery "github.com/nontypeable/financial-tracker/internal/delivery/user"
	accountRepository "github.com/nontypeable/financial-tracker/internal/repository/account"
	userRepository "github.com/nontypeable/financial-tracker/internal/repository/user"
	accountUsecase "github.com/nontypeable/financial-tracker/internal/usecase/account"
	userUsecase "github.com/nontypeable/financial-tracker/internal/usecase/user"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type App struct {
	router chi.Router
	server *http.Server
	config *config.ServerConfig
}

func NewApp(cfg *config.Config, pool *pgxpool.Pool) *App {
	app := App{
		router: chi.NewRouter(),
		config: cfg.Server,
	}

	app.setupMiddleware()
	app.setupRoutes(cfg, pool)

	return &app
}

func (app *App) setupMiddleware() {
	app.router.Use(middleware.RequestID)
	app.router.Use(middleware.RealIP)
	app.router.Use(middleware.Logger)
	app.router.Use(middleware.Recoverer)
	app.router.Use(middleware.Timeout(10 * time.Second))
}

func (app *App) setupRoutes(cfg *config.Config, pool *pgxpool.Pool) {
	app.router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("pong"))
		if err != nil {
			log.Printf("failed to write pong response: %v", err)
		}
	})

	tokenManager := auth.NewTokenManager(cfg.TokenManager.AccessSecret, cfg.TokenManager.RefreshSecret, cfg.TokenManager.AccessTTL, cfg.TokenManager.RefreshTTL)
	authMiddleware := customMiddleware.AuthMiddleware(tokenManager)

	userRepository := userRepository.NewRepository(pool)
	userUsecase := userUsecase.NewService(userRepository, tokenManager)
	userHandler := userDelivery.NewHandler(userUsecase)
	userHandler.RegisterRoutes(app.router, authMiddleware)

	accountRepository := accountRepository.NewRepository(pool)
	accountUsecase := accountUsecase.NewService(accountRepository)
	accountHandler := accountDelivery.NewHandler(accountUsecase)
	accountHandler.RegisterRoutes(app.router, authMiddleware)
}

func (app *App) Start() error {
	if err := app.startServer(); err != nil {
		return err
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh

	log.Printf("Received signal: %s", sig)

	ctx, cancel := context.WithTimeout(context.Background(), app.config.ShutdownTimeout)
	defer cancel()

	return app.Stop(ctx)
}

func (app *App) startServer() error {
	h2s := &http2.Server{}
	handler := h2c.NewHandler(app.router, h2s)

	app.server = &http.Server{
		Addr:    app.config.Address,
		Handler: handler,
	}

	go func() {
		log.Printf("Starting HTTP server on %s", app.config.Address)

		if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	return nil
}

func (app *App) Stop(ctx context.Context) error {
	if app.server == nil {
		return nil
	}

	log.Println("Shutting down HTTP server...")

	if err := app.server.Shutdown(ctx); err != nil {
		log.Printf("Graceful shutdown failed: %v", err)

		if closeErr := app.server.Close(); closeErr != nil {
			log.Printf("Forced close failed: %v", closeErr)
			return closeErr
		}

		return err
	}

	log.Println("Server shut down successfully.")
	return nil
}
