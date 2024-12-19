package run

import (
	"errors"

	"go.temporal.io/sdk/workflow"
)

var (
	ErrInvalidMode = errors.New("invalid mode")
)

type Mode string

const (
	ModeBacktest    Mode = "backtest"
	ModeForwardtest Mode = "forwardtest"
	ModeLive        Mode = "live"
)

func (m Mode) String() string {
	return string(m)
}

func (m Mode) Validate() error {
	switch m {
	case ModeBacktest, ModeForwardtest, ModeLive:
		return nil
	default:
		return ErrInvalidMode
	}
}

const contextModeKey = "cryptellation-run-mode"

func WithMode(ctx workflow.Context, mode Mode) workflow.Context {
	return workflow.WithValue(ctx, contextModeKey, mode)
}

func GetMode(ctx workflow.Context) Mode {
	return ctx.Value(contextModeKey).(Mode)
}
