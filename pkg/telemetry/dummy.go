package telemetry

import "context"

type dummy struct{}

func (d dummy) Close(_ context.Context) {}

func (d dummy) Logger(_ context.Context) Logger {
	return dummyLogger{}
}

func (d dummy) CounterInt(_, _, _ string) (Counter, error) {
	return dummyCounter{}, nil
}

func (d dummy) Trace(ctx context.Context, _, _ string) (context.Context, Tracer) {
	return ctx, dummyTracer{}
}

type dummyLogger struct{}

func (dl dummyLogger) Debug(_ string)            {}
func (dl dummyLogger) Debugf(_ string, _ ...any) {}

func (dl dummyLogger) Info(_ string)            {}
func (dl dummyLogger) Infof(_ string, _ ...any) {}

func (dl dummyLogger) Warning(_ string)            {}
func (dl dummyLogger) Warningf(_ string, _ ...any) {}

func (dl dummyLogger) Error(_ string)            {}
func (dl dummyLogger) Errorf(_ string, _ ...any) {}

type dummyCounter struct{}

func (dl dummyCounter) Add(_ context.Context, _ int64) {}

type dummyTracer struct{}

func (dt dummyTracer) End() {}
