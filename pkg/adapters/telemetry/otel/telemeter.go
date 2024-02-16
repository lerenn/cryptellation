package otel

import (
	"context"
	"fmt"
	"os"

	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// Check implementation of telemetry.Telemetry interface.
var _ telemetry.Telemeter = (*Telemeter)(nil)

type Telemeter struct {
	Logs    logs
	Metrics metrics
	Traces  traces
}

func NewTelemeter(ctx context.Context, serviceName string) (Telemeter, error) {
	// Get exporter URL
	otelEndpoint := os.Getenv("OPENTELEMETRY_GRPC_ENDPOINT")

	// Create Logs exporter
	logs, err := newLogs(ctx, serviceName, otelEndpoint)
	if err != nil {
		return Telemeter{}, err
	}

	// Create Metrics exporter
	metrics, err := newMetrics(ctx, serviceName, otelEndpoint)
	if err != nil {
		return Telemeter{}, err
	}

	// Create Trace exporter
	traces, err := newTraces(ctx, serviceName, otelEndpoint)
	if err != nil {
		metrics.close(ctx)
		return Telemeter{}, err
	}

	return Telemeter{
		Logs:    logs,
		Metrics: metrics,
		Traces:  traces,
	}, nil
}

type counterInt struct {
	counter metric.Int64Counter
}

func (c counterInt) Add(ctx context.Context, value int64) {
	c.counter.Add(ctx, value)
}

func (tel Telemeter) CounterInt(meter, name, description string) (telemetry.Counter, error) {
	desc := metric.WithDescription(description)
	c, err := tel.Metrics.provider.Meter("health").Int64Counter("liveness_calls", desc)

	return counterInt{
		counter: c,
	}, err
}

type Logger struct {
	logger *zap.Logger
}

func (l Logger) Debug(content string) {
	l.logger.Debug(content)
}

func (l Logger) Debugf(format string, a ...any) {
	l.logger.Debug(fmt.Sprintf(format, a...))
}

func (l Logger) Info(content string) {
	l.logger.Info(content)
}

func (l Logger) Infof(format string, a ...any) {
	l.logger.Info(fmt.Sprintf(format, a...))
}

func (l Logger) Warning(content string) {
	l.logger.Warn(content)
}

func (l Logger) Warningf(format string, a ...any) {
	l.logger.Warn(fmt.Sprintf(format, a...))
}

func (l Logger) Error(content string) {
	l.logger.Error(content)
}

func (l Logger) Errorf(format string, a ...any) {
	l.logger.Error(fmt.Sprintf(format, a...))
}

func (tel Telemeter) Logger(ctx context.Context) telemetry.Logger {
	return Logger{
		logger: tel.Logs.logger,
	}
}

type Span struct {
	span trace.Span
}

func (s Span) End() {
	s.span.End()
}

func (tel Telemeter) Trace(ctx context.Context, tracer, name string) (context.Context, telemetry.Tracer) {
	ctx, span := tel.Traces.provider.Tracer(tracer).Start(ctx, name)
	return ctx, Span{
		span: span,
	}
}

func (tel Telemeter) Close(ctx context.Context) {
	tel.Traces.close(ctx)
	tel.Metrics.close(ctx)
	tel.Logs.close(ctx)
}
