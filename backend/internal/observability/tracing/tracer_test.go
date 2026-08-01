package tracing_test

import (
	"context"
	"errors"
	"testing"

	"github.com/cloudstorex/backend/internal/observability/tracing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
)

func TestTracer_InitAndSpans(t *testing.T) {
	shutdown, err := tracing.InitTracer("cloudstorex-test", "test")
	require.NoError(t, err)
	defer func() {
		_ = shutdown(context.Background())
	}()

	// Test root span creation
	ctx, rootSpan := tracing.StartSpan(context.Background(), "test.root")
	defer rootSpan.End()

	traceID := tracing.ExtractTraceID(ctx)
	spanID := tracing.ExtractSpanID(ctx)
	assert.NotEmpty(t, traceID, "Trace ID should be generated")
	assert.NotEmpty(t, spanID, "Span ID should be generated")

	// Test child span creation (Refinement 1: parent/child relationships)
	childCtx, childSpan := tracing.StartChildSpan(ctx, "test.child")
	defer childSpan.End()

	childTraceID := tracing.ExtractTraceID(childCtx)
	childSpanID := tracing.ExtractSpanID(childCtx)
	assert.Equal(t, traceID, childTraceID, "Child span must inherit parent Trace ID")
	assert.NotEqual(t, spanID, childSpanID, "Child span must have distinct Span ID")

	// Test error recording
	errTest := errors.New("simulated error")
	tracing.RecordError(childSpan, errTest)

	// Test nil-safety
	tracing.RecordError(nil, errTest)
	assert.Empty(t, tracing.ExtractTraceID(nil))
	assert.Empty(t, tracing.ExtractSpanID(nil))

	// Test span attributes
	childSpan.SetAttributes(
		attribute.String("workspace.id", "ws-123"),
		attribute.String("provider", "aws-s3"),
	)
}
