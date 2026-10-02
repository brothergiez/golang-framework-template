package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lion/go-framework/internal/infrastructure/cache"
)

type HealthHandler struct {
	pingDB  func(ctx context.Context) error
	cache   cache.Cache
	cacheOn bool
}

func NewHealthHandler(pingDB func(ctx context.Context) error, cacheStore cache.Cache, cacheEnabled bool) *HealthHandler {
	return &HealthHandler{
		pingDB:  pingDB,
		cache:   cacheStore,
		cacheOn: cacheEnabled,
	}
}

func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *HealthHandler) Live(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *HealthHandler) Ready(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	if err := h.pingDB(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not_ready",
			"error":  "database unavailable",
		})
		return
	}
	if h.cacheOn {
		if err := h.cache.Ping(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not_ready",
				"error":  "cache unavailable",
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}
