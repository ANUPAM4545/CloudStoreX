package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/cloudstorex/backend/internal/observability/metrics"
	"github.com/gin-gonic/gin"
)

// PrometheusMiddleware collects HTTP RED (Rate, Errors, Duration) signals for incoming requests.
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		route := c.FullPath()
		if route == "" {
			route = c.Request.URL.Path
		}
		method := c.Request.Method

		// Gauge active requests
		activeGauge := metrics.HTTPActiveRequests.WithLabelValues(method, route)
		activeGauge.Inc()
		defer activeGauge.Dec()

		c.Next()

		duration := time.Since(start).Seconds()
		status := c.Writer.Status()
		statusStr := strconv.Itoa(status)

		// Record RED request rate and latency
		metrics.HTTPRequestsTotal.WithLabelValues(method, route, statusStr).Inc()
		metrics.HTTPRequestDurationSeconds.WithLabelValues(method, route, statusStr).Observe(duration)

		// Record HTTP error count if 4xx or 5xx
		if status >= http.StatusBadRequest {
			metrics.HTTPErrorRate.WithLabelValues(method, route, statusStr).Inc()
		}
	}
}
