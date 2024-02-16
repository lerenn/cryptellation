package telemetry

import "context"

type dummy struct{}

func (d dummy) Close(ctx context.Context) {}

func (d dummy) Logger(ctx context.Context) Logger {
	return dummyLogger{}
}

func (d dummy) CounterInt(meter, name, description string) (Counter, error) {
	return dummyCounter{}, nil
}

func (d dummy) Trace(ctx context.Context, tracer, name string) (context.Context, Tracer) {
	return ctx, dummyTracer{}
}

type dummyLogger struct{}

func (dl dummyLogger) Debug(text string)              {}
func (dl dummyLogger) Debugf(format string, a ...any) {}

func (dl dummyLogger) Info(text string)              {}
func (dl dummyLogger) Infof(format string, a ...any) {}

func (dl dummyLogger) Warning(text string)              {}
func (dl dummyLogger) Warningf(format string, a ...any) {}

func (dl dummyLogger) Error(text string)              {}
func (dl dummyLogger) Errorf(format string, a ...any) {}

type dummyCounter struct{}

func (dl dummyCounter) Add(ctx context.Context, value int64) {}

type dummyTracer struct{}

func (dt dummyTracer) End() {}
