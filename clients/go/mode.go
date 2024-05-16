package client

import "errors"

var (
	ErrUnsupportedMode = errors.New("unsupported mode")
)

type Mode string

const (
	ModeIsBacktest Mode = "backtest"
	ModeIsForward  Mode = "forward"
	ModeIsLive     Mode = "live"
)

func (m Mode) String() string {
	return string(m)
}

func (m Mode) IsValid() bool {
	switch m {
	case ModeIsBacktest, ModeIsForward, ModeIsLive:
		return true
	}
	return false
}
