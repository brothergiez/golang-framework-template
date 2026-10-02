package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lion/go-framework/internal/app/bootstrap"
	"github.com/lion/go-framework/internal/config"
)

func main() {
	if os.Getenv("SERVICE_TYPE") == "" {
		_ = os.Setenv("SERVICE_TYPE", config.ServiceTypeWorker)
	}

	app, err := bootstrap.New()
	if err != nil {
		log.Fatalf("bootstrap failed: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx); err != nil {
		log.Fatalf("application stopped: %v", err)
	}
}
