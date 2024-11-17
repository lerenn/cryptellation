package otel

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

type traces struct {
	exporter *otlptrace.Exporter
	provider *trace.TracerProvider
}

func newTraces(ctx context.Context, serviceName, url string) (traces, error) {
	// Create exporter
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint(url), otlptracegrpc.WithInsecure())
	if err != nil {
		return traces{}, err
	}

	// Create resource
	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
	)

	// Create provider
	provider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource),
		// Set the sampling rate based on the parent span to 60%
		trace.WithSampler(trace.ParentBased(trace.TraceIDRatioBased(0.6))),
	)

	// Set opentelemetry traces provider globally
	// otel.SetTracerProvider(provider)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{}, // W3C Trace Context format; https://www.w3.org/TR/trace-context/
		),
	)

	return traces{
		exporter: exporter,
		provider: provider,
	}, nil
}

func (t traces) close(ctx context.Context) {
	_ = t.provider.Shutdown(ctx)
	_ = t.exporter.Shutdown(ctx)
}
