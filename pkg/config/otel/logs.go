package otel

import (
	"context"

	otellogs "github.com/agoda-com/opentelemetry-logs-go"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs/otlplogsgrpc"
	sdk "github.com/agoda-com/opentelemetry-logs-go/sdk/logs"
	"github.com/agoda-com/otelzap"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.uber.org/zap"
)

type Logs struct {
	clientCleanUp func(ctx context.Context) error
	exporter      *otlplogs.Exporter
	provider      *sdk.LoggerProvider
	logger        *zap.Logger
	undo          func()
}

func NewLogs(ctx context.Context, serviceName, url string) (Logs, error) {
	context.Background()

	// Create exporter
	client := otlplogsgrpc.NewClient(otlplogsgrpc.WithEndpoint(url), otlplogsgrpc.WithInsecure())
	exporter, err := otlplogs.NewExporter(ctx, otlplogs.WithClient(client))
	if err != nil {
		return Logs{}, err
	}

	// Create resource
	resource := resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String(serviceName))

	// Create provider
	provider := sdk.NewLoggerProvider(
		sdk.WithBatcher(exporter),
		sdk.WithResource(resource),
	)

	// Set opentelemetry logger provider globally
	otellogs.SetLoggerProvider(provider)

	// Create a new logger
	logger := zap.New(otelzap.NewOtelCore(provider))

	return Logs{
		clientCleanUp: client.Stop,
		exporter:      exporter,
		provider:      provider,
		logger:        logger,
		undo:          zap.ReplaceGlobals(logger),
	}, nil
}

func (l Logs) Close(ctx context.Context) {
	l.undo()
	l.logger.Sync()
	l.provider.Shutdown(ctx)
	l.exporter.Shutdown(ctx)
	_ = l.clientCleanUp(ctx)
}
