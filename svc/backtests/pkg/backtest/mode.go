package backtest

import "errors"

var (
	ErrInvalidMode = errors.New("invalid mode")
)

type Mode string

const (
	ModeIsFullOHLC  Mode = "full_ohlc"
	ModeIsCloseOHLC Mode = "close_ohlc"
)

func (m Mode) Validate() error {
	switch m {
	case ModeIsFullOHLC, ModeIsCloseOHLC:
		return nil
	default:
		return ErrInvalidMode
	}
}

func (m Mode) String() string {
	return string(m)
}

func (m Mode) Opt() *Mode {
	return &m
}
