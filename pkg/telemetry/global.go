package telemetry

import "context"

var (
	globalTelemeter Telemeter = dummy{}
)

// Set sets the global telemeter.
func Set(t Telemeter) {
	globalTelemeter = t
}

// L returns a logger from the global telemeter.
func L(ctx context.Context) Logger {
	return globalTelemeter.Logger(ctx)
}

// T returns a tracer from the global telemeter.
func T(ctx context.Context, tracer, name string) (context.Context, Tracer) {
	return globalTelemeter.Trace(ctx, tracer, name)
}

// Ci returns an int counter from the global telemeter.
func Ci(meter, name, description string) (Counter, error) {
	return globalTelemeter.CounterInt(meter, name, description)
}

// Close closes the global telemeter.
func Close(ctx context.Context) {
	globalTelemeter.Close(ctx)
}
