package middleware

import (
	"log/slog"
	"time"

	"github.com/cloudstorex/backend/internal/shared/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Generate Request ID
		requestID := uuid.New().String()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		userID, _ := c.Get("user_id")
		traceID, _ := c.Get("trace_id")
		spanID, _ := c.Get("span_id")
		workspaceID, _ := c.Get("workspace_id")

		logger.Log.Info("HTTP Request",
			slog.String("request_id", requestID),
			slog.Any("trace_id", traceID),
			slog.Any("span_id", spanID),
			slog.Any("workspace_id", workspaceID),
			slog.Any("user_id", userID),
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.String("query", query),
			slog.String("ip", c.ClientIP()),
			slog.String("user_agent", c.Request.UserAgent()),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("latency", latency),
		)
	}
}
