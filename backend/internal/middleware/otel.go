package middleware

import (
	"fmt"
	"net/http"

	"github.com/cloudstorex/backend/internal/observability/tracing"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
)

// OpenTelemetryMiddleware creates a Gin middleware that starts an OpenTelemetry span
// for each incoming HTTP request, propagating W3C trace context and setting child context.
func OpenTelemetryMiddleware(serviceName string) gin.HandlerFunc {
	propagator := otel.GetTextMapPropagator()

	return func(c *gin.Context) {
		// Extract incoming trace context from HTTP request headers
		ctx := propagator.Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))

		route := c.FullPath()
		if route == "" {
			route = c.Request.URL.Path
		}
		spanName := fmt.Sprintf("HTTP %s %s", c.Request.Method, route)

		ctx, span := tracing.StartSpan(ctx, spanName, trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		traceID := tracing.ExtractTraceID(ctx)
		spanID := tracing.ExtractSpanID(ctx)

		// Set trace context in Gin Context and HTTP response headers
		c.Set("trace_id", traceID)
		c.Set("span_id", spanID)
		c.Header("X-Trace-ID", traceID)
		c.Header("X-Span-ID", spanID)

		// Set semantic HTTP attributes
		span.SetAttributes(
			semconv.HTTPMethod(c.Request.Method),
			semconv.HTTPRoute(route),
			semconv.HTTPTarget(c.Request.URL.Path),
			attribute.String("http.user_agent", c.Request.UserAgent()),
			attribute.String("http.client_ip", c.ClientIP()),
		)

		// Replace request context so downstream handlers inherit the span (Refinement 1)
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		status := c.Writer.Status()
		span.SetAttributes(semconv.HTTPStatusCode(status))

		if status >= http.StatusInternalServerError {
			tracing.RecordError(span, fmt.Errorf("HTTP %d error", status))
		}
	}
}
