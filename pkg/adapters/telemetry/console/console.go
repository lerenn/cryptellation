package console

import (
	"context"
	"fmt"

	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
)

type Telemeter struct{}

func (t Telemeter) Close(ctx context.Context) {}

func (t Telemeter) Logger(ctx context.Context) telemetry.Logger {
	return consoleLogger{}
}

func (t Telemeter) CounterInt(meter, name, description string) (telemetry.Counter, error) {
	return consoleCounter{}, nil
}

func (t Telemeter) Trace(ctx context.Context, tracer, name string) (context.Context, telemetry.Tracer) {
	return ctx, consoleTracer{}
}

type consoleLogger struct{}

func (cl consoleLogger) Debug(text string) {
	fmt.Println("DEBUG:", text)
}

func (cl consoleLogger) Debugf(format string, a ...any) {
	fmt.Printf("DEBUG: "+format, a...)
	fmt.Println()
}

func (cl consoleLogger) Info(text string) {
	fmt.Println("INFO:", text)
}

func (cl consoleLogger) Infof(format string, a ...any) {
	fmt.Printf("INFO: "+format, a...)
	fmt.Println()
}

func (cl consoleLogger) Warning(text string) {
	fmt.Println("WARNING:", text)
}

func (cl consoleLogger) Warningf(format string, a ...any) {
	fmt.Printf("WARNING: "+format, a...)
	fmt.Println()
}

func (cl consoleLogger) Error(text string) {
	fmt.Println("ERROR:", text)
}

func (cl consoleLogger) Errorf(format string, a ...any) {
	fmt.Printf("ERROR: "+format, a...)
	fmt.Println()
}

type consoleCounter struct{}

func (cl consoleCounter) Add(ctx context.Context, value int64) {}

type consoleTracer struct{}

func (ct consoleTracer) End() {}
