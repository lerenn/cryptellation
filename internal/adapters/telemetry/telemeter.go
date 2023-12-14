package telemetry

import "context"

type Telemeter interface {
	Close(ctx context.Context)

	// Logs
	Logger(ctx context.Context) Logger

	// Metrics
	CounterInt(meter, name, description string) (Counter, error)

	// Traces
	Trace(ctx context.Context, tracer, name string) (context.Context, Tracer)
}

type Logger interface {
	Debug(text string)
	Error(text string)
}

type Tracer interface {
	End()
}

type Counter interface {
	Add(ctx context.Context, value int64)
}
