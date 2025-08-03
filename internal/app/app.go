package app

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nontypeable/financial-tracker/internal/auth"
	"github.com/nontypeable/financial-tracker/internal/config"
	userDelivery "github.com/nontypeable/financial-tracker/internal/delivery/user"
	userRepository "github.com/nontypeable/financial-tracker/internal/repository/user"
	userUsecase "github.com/nontypeable/financial-tracker/internal/usecase/user"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type App struct {
	router chi.Router
	server *http.Server
	config *config.ServerConfig
}

func NewApp(cfg *config.Config, db *sql.DB) *App {
	app := App{
		router: chi.NewRouter(),
		config: cfg.Server,
	}

	app.setupMiddleware()
	app.setupRoutes(cfg, db)

	return &app
}

func (app *App) setupMiddleware() {
	app.router.Use(middleware.RequestID)
	app.router.Use(middleware.RealIP)
	app.router.Use(middleware.Logger)
	app.router.Use(middleware.Recoverer)
	app.router.Use(middleware.Timeout(10 * time.Second))
}

func (app *App) setupRoutes(cfg *config.Config, db *sql.DB) {
	app.router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	tokenManager := auth.NewTokenManager(cfg.TokenManager.AccessSecret, cfg.TokenManager.RefreshSecret, cfg.TokenManager.AccessTTL, cfg.TokenManager.RefreshTTL)

	userRepository := userRepository.NewRepository(db)
	userUsecase := userUsecase.NewService(userRepository, tokenManager)
	userHandler := userDelivery.NewHandler(userUsecase)
	userHandler.RegisterRoutes(app.router)
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
