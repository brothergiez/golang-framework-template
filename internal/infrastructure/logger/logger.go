package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/lion/go-framework/internal/config"
)

type ctxKey string

const (
	RequestIDKey     ctxKey = "request_id"
	CorrelationIDKey ctxKey = "correlation_id"
)

func New(cfg config.LoggerConfig) *slog.Logger {
	level := parseLevel(cfg.Level)
	opts := &slog.HandlerOptions{Level: level}

	var handler slog.Handler
	if strings.ToLower(cfg.Format) == "text" {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}
	return slog.New(handler)
}

func WithContext(ctx context.Context, log *slog.Logger) *slog.Logger {
	if log == nil {
		log = slog.Default()
	}
	if v, ok := ctx.Value(RequestIDKey).(string); ok && v != "" {
		log = log.With("request_id", v)
	}
	if v, ok := ctx.Value(CorrelationIDKey).(string); ok && v != "" {
		log = log.With("correlation_id", v)
	}
	return log
}

func ContextWithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

func ContextWithCorrelationID(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, CorrelationIDKey, correlationID)
}

func RequestIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(RequestIDKey).(string); ok {
		return v
	}
	return ""
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
