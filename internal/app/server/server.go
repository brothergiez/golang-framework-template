package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lion/go-framework/internal/config"
)

type Server struct {
	httpServer *http.Server
	log        *slog.Logger
}

func New(cfg *config.Config, engine *gin.Engine, log *slog.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         cfg.HTTPAddr(),
			Handler:      engine,
			ReadTimeout:  cfg.HTTP.ReadTimeout,
			WriteTimeout: cfg.HTTP.WriteTimeout,
		},
		log: log,
	}
}

func (s *Server) Start() error {
	s.log.Info("api server starting", "addr", s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("listen: %w", err)
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("api server shutting down")
	return s.httpServer.Shutdown(ctx)
}
