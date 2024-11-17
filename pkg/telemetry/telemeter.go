package telemetry

import "context"

// Telemeter is the interface for the telemetry system.
type Telemeter interface {
	// Close the telemetry system.
	Close(ctx context.Context)

	// Logs
	Logger(ctx context.Context) Logger

	// Metrics
	CounterInt(meter, name, description string) (Counter, error)

	// Traces
	Trace(ctx context.Context, tracer, name string) (context.Context, Tracer)
}

// Logger is the logging interface for the telemetry system.
type Logger interface {
	// Log a debug message.
	Debug(text string)
	// Log a formatted debug message.
	Debugf(format string, a ...any)

	// Log an info message.
	Info(text string)
	// Log a formatted info message.
	Infof(format string, a ...any)

	// Log a warning message.
	Warning(text string)
	// Log a formatted warning message.
	Warningf(format string, a ...any)

	// Log an error message.
	Error(text string)
	// Log a formatted error message.
	Errorf(format string, a ...any)
}

// Tracer is the tracing interface for the telemetry system.
type Tracer interface {
	// Stop the tracing.
	End()
}

// Counter is the interface for a counter.
type Counter interface {
	// Add a value to the counter.
	Add(ctx context.Context, value int64)
}
