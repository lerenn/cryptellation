package telemetry

import "context"

var (
	globalTelemeter Telemeter = dummy{}
)

func Set(t Telemeter) {
	globalTelemeter = t
}

func L(ctx context.Context) Logger {
	return globalTelemeter.Logger(ctx)
}

func T(ctx context.Context, tracer, name string) (context.Context, Tracer) {
	return globalTelemeter.Trace(ctx, tracer, name)
}

func CI(meter, name, description string) (Counter, error) {
	return globalTelemeter.CounterInt(meter, name, description)
}
