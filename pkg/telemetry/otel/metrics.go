package otel

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

type Metrics struct {
	exporter *otlpmetricgrpc.Exporter
	provider *sdkmetric.MeterProvider
}

func NewMetrics(ctx context.Context, serviceName, url string) (Metrics, error) {
	// Create an exporter
	exporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithEndpoint(url), otlpmetricgrpc.WithInsecure())
	if err != nil {
		return Metrics{}, err
	}

	// Create resource
	resource := resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String(serviceName))

	// Create provider
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(resource),
		sdkmetric.WithReader(
			// collects and exports metric data every 30 seconds.
			sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(30*time.Second)),
		),
	)

	// Set opentelemetry metrics provider globally
	otel.SetMeterProvider(provider)

	return Metrics{
		exporter: exporter,
		provider: provider,
	}, nil
}

func (m Metrics) Close(ctx context.Context) {
	m.provider.Shutdown(ctx)
	m.exporter.Shutdown(ctx)
}
