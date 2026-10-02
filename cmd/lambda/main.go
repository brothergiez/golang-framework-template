package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lion/go-framework/internal/app/bootstrap"
	"github.com/lion/go-framework/internal/config"
)

func main() {
	if os.Getenv("SERVICE_TYPE") == "" {
		_ = os.Setenv("SERVICE_TYPE", config.ServiceTypeLambda)
	}

	app, err := bootstrap.New()
	if err != nil {
		log.Fatalf("bootstrap failed: %v", err)
	}
	defer func() { _ = app.Close() }()

	event := []byte(os.Getenv("LAMBDA_EVENT"))
	if len(event) == 0 {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			event, err = io.ReadAll(os.Stdin)
			if err != nil {
				log.Fatalf("read stdin: %v", err)
			}
		}
	}
	if len(event) == 0 {
		event, _ = json.Marshal(map[string]string{
			"action": "user.create",
			"name":   "Lambda User",
			"email":  "lambda@example.com",
		})
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	out, err := app.Lambda.Invoke(ctx, event)
	if err != nil {
		log.Fatalf("invoke failed: %v", err)
	}
	_, _ = os.Stdout.Write(out)
	_, _ = os.Stdout.Write([]byte("\n"))
}
