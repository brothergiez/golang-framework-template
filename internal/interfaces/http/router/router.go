package router

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/lion/go-framework/internal/config"
	"github.com/lion/go-framework/internal/infrastructure/metrics"
	"github.com/lion/go-framework/internal/interfaces/http/handler"
	"github.com/lion/go-framework/internal/interfaces/http/middleware"
)

type Dependencies struct {
	Config  *config.Config
	Logger  *slog.Logger
	Users   *handler.UserHandler
	Health  *handler.HealthHandler
}

func New(deps Dependencies) *gin.Engine {
	if deps.Config.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(
		middleware.RequestID(),
		middleware.Recovery(deps.Logger),
		middleware.AccessLog(deps.Logger),
		middleware.Metrics(),
		middleware.SecurityHeaders(),
		middleware.BodyLimit(deps.Config.HTTP.BodyLimit),
	)

	r.GET("/health", deps.Health.Health)
	r.GET("/health/live", deps.Health.Live)
	r.GET("/health/ready", deps.Health.Ready)
	r.GET("/metrics", gin.WrapH(metrics.Handler()))

	v1 := r.Group("/api/v1")
	{
		users := v1.Group("/users")
		users.POST("", deps.Users.Create)
		users.GET("", deps.Users.List)
		users.GET("/:id", deps.Users.GetByID)
		users.PUT("/:id", deps.Users.Update)
		users.DELETE("/:id", deps.Users.Delete)
	}

	return r
}
