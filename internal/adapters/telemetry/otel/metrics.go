package otel

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

type metrics struct {
	exporter *otlpmetricgrpc.Exporter
	provider *sdkmetric.MeterProvider
}

func newMetrics(ctx context.Context, serviceName, url string) (metrics, error) {
	// Create an exporter
	exporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithEndpoint(url), otlpmetricgrpc.WithInsecure())
	if err != nil {
		return metrics{}, err
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
	// otel.SetMeterProvider(provider)

	return metrics{
		exporter: exporter,
		provider: provider,
	}, nil
}

func (m metrics) close(ctx context.Context) {
	_ = m.provider.Shutdown(ctx)
	_ = m.exporter.Shutdown(ctx)
}
