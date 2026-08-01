package tracing

import (
	"context"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

const TracerName = "cloudstorex.backend"

// InitTracer initializes the OpenTelemetry distributed tracing provider.
// If OTEL_EXPORTER_OTLP_ENDPOINT is configured, it exports traces via OTLP HTTP.
// Otherwise, it sets up a no-op trace provider to prevent panics during dev/testing.
func InitTracer(serviceName, env string) (func(context.Context) error, error) {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")

	var tp trace.TracerProvider
	if endpoint != "" {
		exporter, err := otlptracehttp.New(context.Background(),
			otlptracehttp.WithEndpoint(endpoint),
			otlptracehttp.WithInsecure(),
		)
		if err != nil {
			return nil, err
		}

		res, err := resource.New(context.Background(),
			resource.WithAttributes(
				semconv.ServiceName(serviceName),
				attribute.String("deployment.environment", env),
			),
		)
		if err != nil {
			return nil, err
		}

		sdkTP := sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(res),
		)
		tp = sdkTP
		otel.SetTracerProvider(tp)
	} else {
		// Use default SDK provider without exporter for unit tests and local dev
		res := resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			attribute.String("deployment.environment", env),
		)
		sdkTP := sdktrace.NewTracerProvider(
			sdktrace.WithResource(res),
		)
		tp = sdkTP
		otel.SetTracerProvider(tp)
	}

	// Set global W3C trace context and baggage propagators
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return func(ctx context.Context) error {
		if sdkTP, ok := tp.(*sdktrace.TracerProvider); ok {
			return sdkTP.Shutdown(ctx)
		}
		return nil
	}, nil
}

// StartSpan starts a new root or child span using the global tracer.
// It maintains proper parent/child relationships from the provided context (Refinement 1).
func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	if ctx == nil {
		ctx = context.Background()
	}
	return otel.Tracer(TracerName).Start(ctx, name, opts...)
}

// StartChildSpan explicitly starts a child span inheriting from parent context.
func StartChildSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return StartSpan(ctx, name, opts...)
}

// RecordError records an error on the span and sets span status to Error.
func RecordError(span trace.Span, err error) {
	if span == nil || err == nil {
		return
	}
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

// ExtractTraceID returns the Trace ID hex string from context, or empty if none exists.
func ExtractTraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	spanCtx := trace.SpanFromContext(ctx).SpanContext()
	if !spanCtx.IsValid() {
		return ""
	}
	return spanCtx.TraceID().String()
}

// ExtractSpanID returns the Span ID hex string from context, or empty if none exists.
func ExtractSpanID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	spanCtx := trace.SpanFromContext(ctx).SpanContext()
	if !spanCtx.IsValid() {
		return ""
	}
	return spanCtx.SpanID().String()
}

// NoopSpan returns a noop span for testing when tracing is disabled.
func NoopSpan() trace.Span {
	return noop.Span{}
}
