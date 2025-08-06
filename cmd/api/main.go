package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nontypeable/financial-tracker/internal/app"
	"github.com/nontypeable/financial-tracker/internal/config"
	"github.com/nontypeable/financial-tracker/internal/repository"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.LoadConfig(os.Getenv("CONFIG_PATH"))
	if err != nil {
		log.Fatal(err.Error())
	}

	pool, err := repository.NewPostgres(ctx, cfg.Database)
	if err != nil {
		log.Fatal(err.Error())
	}

	application := app.NewApp(cfg, pool)
	if err := application.Start(); err != nil {
		log.Fatalf("application error: %v", err)
	}
}
