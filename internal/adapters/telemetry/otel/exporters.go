package otel

import (
	"context"
	"os"
)

type Exporters struct {
	Logs    Logs
	Metrics Metrics
	Traces  Traces
}

func NewExporters(ctx context.Context, serviceName string) (Exporters, error) {
	// Get exporter URL
	otelEndpoint := os.Getenv("OPENTELEMETRY_GRPC_ENDPOINT")

	// Create Logs exporter
	logs, err := NewLogs(ctx, serviceName, otelEndpoint)
	if err != nil {
		return Exporters{}, err
	}

	// Create Metrics exporter
	metrics, err := NewMetrics(ctx, serviceName, otelEndpoint)
	if err != nil {
		return Exporters{}, err
	}

	// Create Trace exporter
	traces, err := NewTraces(ctx, serviceName, otelEndpoint)
	if err != nil {
		metrics.Close(ctx)
		return Exporters{}, err
	}

	return Exporters{
		Logs:    logs,
		Metrics: metrics,
		Traces:  traces,
	}, nil
}

func (exp Exporters) Close(ctx context.Context) {
	exp.Traces.Close(ctx)
	exp.Metrics.Close(ctx)
	exp.Logs.Close(ctx)
}
