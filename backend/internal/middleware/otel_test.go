package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudstorex/backend/internal/middleware"
	"github.com/cloudstorex/backend/internal/observability/tracing"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenTelemetryMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	shutdown, err := tracing.InitTracer("cloudstorex-test", "test")
	require.NoError(t, err)
	defer func() {
		_ = shutdown(context.Background())
	}()

	r := gin.New()
	r.Use(middleware.OpenTelemetryMiddleware("cloudstorex-test"))
	r.GET("/api/v1/test", func(c *gin.Context) {
		traceID, exists := c.Get("trace_id")
		assert.True(t, exists, "trace_id should be stored in gin context")
		assert.NotEmpty(t, traceID, "trace_id should not be empty")

		// Verify child span creation inside handler (Refinement 1)
		childCtx, childSpan := tracing.StartChildSpan(c.Request.Context(), "handler.child")
		defer childSpan.End()

		childTraceID := tracing.ExtractTraceID(childCtx)
		assert.Equal(t, traceID, childTraceID, "Child span must inherit parent request trace ID")

		c.String(http.StatusOK, "OK")
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Header().Get("X-Trace-ID"), "X-Trace-ID response header should be set")
	assert.NotEmpty(t, w.Header().Get("X-Span-ID"), "X-Span-ID response header should be set")
}
