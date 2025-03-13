package otel

import (
	"context"
	"fmt"
	"os"

	"github.com/lerenn/cryptellation/v1/pkg/config"
	"github.com/lerenn/cryptellation/v1/pkg/telemetry"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// Check implementation of telemetry.Telemetry interface.
var _ telemetry.Telemeter = (*telemeter)(nil)

type telemeter struct {
	Logs    logs
	Metrics metrics
	Traces  traces
}

// NewTelemeter creates a new telemeter based on OpenTelemetry.
func NewTelemeter(ctx context.Context, serviceName string) (telemetry.Telemeter, error) {
	// Get exporter URL
	otelEndpoint := os.Getenv(config.EnvOpentelemetryGrpcEndpoint)

	// Create Logs exporter
	logs, err := newLogs(ctx, serviceName, otelEndpoint)
	if err != nil {
		return telemeter{}, err
	}

	// Create Metrics exporter
	metrics, err := newMetrics(ctx, serviceName, otelEndpoint)
	if err != nil {
		return telemeter{}, err
	}

	// Create Trace exporter
	traces, err := newTraces(ctx, serviceName, otelEndpoint)
	if err != nil {
		metrics.close(ctx)
		return telemeter{}, err
	}

	return telemeter{
		Logs:    logs,
		Metrics: metrics,
		Traces:  traces,
	}, nil
}

type counterInt struct {
	counter metric.Int64Counter
}

// Add adds a value to the counter.
func (c counterInt) Add(ctx context.Context, value int64) {
	c.counter.Add(ctx, value)
}

// CounterInt creates a new integer counter.
func (tel telemeter) CounterInt(meter, name, description string) (telemetry.Counter, error) {
	desc := metric.WithDescription(description)
	c, err := tel.Metrics.provider.Meter(meter).Int64Counter(name, desc)

	return counterInt{
		counter: c,
	}, err
}

// Logger is the logger struct for OpenTelemetry.
type Logger struct {
	logger *zap.Logger
}

// Debug logs a debug message.
func (l Logger) Debug(content string) {
	l.logger.Debug(content)
}

// Debugf logs a formatted debug message.
func (l Logger) Debugf(format string, a ...any) {
	l.logger.Debug(fmt.Sprintf(format, a...))
}

// Info logs an info message.
func (l Logger) Info(content string) {
	l.logger.Info(content)
}

// Infof logs a formatted info message.
func (l Logger) Infof(format string, a ...any) {
	l.logger.Info(fmt.Sprintf(format, a...))
}

// Warning logs a warning message.
func (l Logger) Warning(content string) {
	l.logger.Warn(content)
}

// Warningf logs a formatted warning message.
func (l Logger) Warningf(format string, a ...any) {
	l.logger.Warn(fmt.Sprintf(format, a...))
}

// Error logs an error message.
func (l Logger) Error(content string) {
	l.logger.Error(content)
}

// Errorf logs a formatted error message.
func (l Logger) Errorf(format string, a ...any) {
	l.logger.Error(fmt.Sprintf(format, a...))
}

// Logger returns a new logger.
func (tel telemeter) Logger(_ context.Context) telemetry.Logger {
	return Logger{
		logger: tel.Logs.logger,
	}
}

// Span is the tracing struct for OpenTelemetry.
type Span struct {
	span trace.Span
}

// End stops the tracing.
func (s Span) End() {
	s.span.End()
}

// Trace starts a new trace.
func (tel telemeter) Trace(ctx context.Context, tracer, name string) (context.Context, telemetry.Tracer) {
	ctx, span := tel.Traces.provider.Tracer(tracer).Start(ctx, name)
	return ctx, Span{
		span: span,
	}
}

// Close closes the telemeter.
func (tel telemeter) Close(ctx context.Context) {
	tel.Traces.close(ctx)
	tel.Metrics.close(ctx)
	tel.Logs.close(ctx)
}
