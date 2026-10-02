package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lion/go-framework/internal/infrastructure/logger"
	"github.com/lion/go-framework/internal/infrastructure/metrics"
	"github.com/lion/go-framework/internal/shared/utils"
	"log/slog"
)

const (
	HeaderRequestID     = "X-Request-ID"
	HeaderCorrelationID = "X-Correlation-ID"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader(HeaderRequestID)
		if rid == "" {
			rid = utils.NewID()
		}
		cid := c.GetHeader(HeaderCorrelationID)
		if cid == "" {
			cid = rid
		}
		ctx := logger.ContextWithRequestID(c.Request.Context(), rid)
		ctx = logger.ContextWithCorrelationID(ctx, cid)
		c.Request = c.Request.WithContext(ctx)
		c.Writer.Header().Set(HeaderRequestID, rid)
		c.Writer.Header().Set(HeaderCorrelationID, cid)
		c.Set(string(logger.RequestIDKey), rid)
		c.Next()
	}
}

func Recovery(log *slog.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		logger.WithContext(c.Request.Context(), log).Error("panic recovered", "panic", recovered)
		c.AbortWithStatusJSON(500, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "internal server error",
			},
		})
	})
}

func AccessLog(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.WithContext(c.Request.Context(), log).Info("http_request",
			"method", c.Request.Method,
			"path", c.FullPath(),
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
			"client_ip", c.ClientIP(),
		)
	}
}

func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		status := strconv.Itoa(c.Writer.Status())
		metrics.HTTPRequests.WithLabelValues(c.Request.Method, path, status).Inc()
		metrics.HTTPDuration.WithLabelValues(c.Request.Method, path).Observe(time.Since(start).Seconds())
	}
}

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("Referrer-Policy", "no-referrer")
		c.Next()
	}
}

func BodyLimit(limit int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = httpMaxBytesReader(c.Writer, c.Request.Body, limit)
		c.Next()
	}
}
