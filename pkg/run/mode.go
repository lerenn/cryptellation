package run

import (
	"errors"

	"go.temporal.io/sdk/workflow"
)

var (
	// ErrInvalidMode is returned when the mode is invalid.
	ErrInvalidMode = errors.New("invalid mode")
)

// Mode is the mode in which a run is executed.
type Mode string

const (
	// ModeBacktest is the backtest mode.
	ModeBacktest Mode = "backtest"
	// ModeForwardtest is the forwardtest mode.
	ModeForwardtest Mode = "forwardtest"
	// ModeLive is the live mode.
	ModeLive Mode = "live"
)

// String returns the string representation of the mode.
func (m Mode) String() string {
	return string(m)
}

// Validate validates the mode.
func (m Mode) Validate() error {
	switch m {
	case ModeBacktest, ModeForwardtest, ModeLive:
		return nil
	default:
		return ErrInvalidMode
	}
}

const contextModeKey = "cryptellation-run-mode"

// WithMode sets the mode in the context.
func WithMode(ctx workflow.Context, mode Mode) workflow.Context {
	return workflow.WithValue(ctx, contextModeKey, mode)
}

// GetMode returns the mode from the context.
func GetMode(ctx workflow.Context) Mode {
	return ctx.Value(contextModeKey).(Mode)
}
