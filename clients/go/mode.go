package client

import "errors"

var (
	ErrUnsupportedMode = errors.New("unsupported mode")
)

type Mode string

const (
	ModeIsBacktest    Mode = "backtest"
	ModeIsForwardTest Mode = "forwardtest"
	ModeIsLive        Mode = "live"
)

func (m Mode) String() string {
	return string(m)
}

func (m Mode) IsValid() bool {
	switch m {
	case ModeIsBacktest, ModeIsForwardTest, ModeIsLive:
		return true
	}
	return false
}
