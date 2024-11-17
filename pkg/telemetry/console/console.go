package console

import (
	"context"
	"fmt"
	"time"

	"github.com/lerenn/cryptellation/v1/pkg/telemetry"
)

// Telemeter is a console telemeter.
type Telemeter struct{}

// Close closes the console telemeter.
func (t Telemeter) Close(ctx context.Context) {}

// Logger returns a console logger.
func (t Telemeter) Logger(ctx context.Context) telemetry.Logger {
	return consoleLogger{}
}

// CounterInt returns a console integer counter.
func (t Telemeter) CounterInt(meter, name, description string) (telemetry.Counter, error) {
	return consoleIntCounter{}, nil
}

// Trace returns a console tracer.
func (t Telemeter) Trace(ctx context.Context, tracer, name string) (context.Context, telemetry.Tracer) {
	return ctx, consoleTracer{}
}

type consoleLogger struct{}

func (cl consoleLogger) print(level, text string) {
	now := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%s %-5s %s\n", now, level, text)
}

// Debug logs a debug message.
func (cl consoleLogger) Debug(text string) {
	cl.print("DEBUG", text)
}

// Debugf logs a formatted debug message.
func (cl consoleLogger) Debugf(format string, a ...any) {
	cl.print("DEBUG", fmt.Sprintf(format, a...))
}

// Info logs an info message.
func (cl consoleLogger) Info(text string) {
	cl.print("INFO", text)
}

// Infof logs a formatted info message.
func (cl consoleLogger) Infof(format string, a ...any) {
	cl.print("INFO", fmt.Sprintf(format, a...))
}

// Warning logs a warning message.
func (cl consoleLogger) Warning(text string) {
	cl.print("WARN", text)
}

// Warningf logs a formatted warning message.
func (cl consoleLogger) Warningf(format string, a ...any) {
	cl.print("WARN", fmt.Sprintf(format, a...))
}

// Error logs an error message.
func (cl consoleLogger) Error(text string) {
	cl.print("ERR", text)
}

// Errorf logs a formatted error message.
func (cl consoleLogger) Errorf(format string, a ...any) {
	cl.print("ERR", fmt.Sprintf(format, a...))
}

type consoleIntCounter struct{}

// Add adds a value to the console integer counter.
func (cl consoleIntCounter) Add(ctx context.Context, value int64) {}

type consoleTracer struct{}

// End ends the console tracer.
func (ct consoleTracer) End() {}
