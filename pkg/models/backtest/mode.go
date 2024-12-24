package backtest

import "errors"

var (
	// ErrInvalidMode is returned when the mode is invalid.
	ErrInvalidMode = errors.New("invalid mode")
)

// Mode is the mode of the backtest.
type Mode string

const (
	// ModeIsFullOHLC is the mode where the backtest uses full OHLC data.
	ModeIsFullOHLC Mode = "full_ohlc"
	// ModeIsCloseOHLC is the mode where the backtest uses close OHLC data.
	ModeIsCloseOHLC Mode = "close_ohlc"
)

// Validate will validate the mode.
func (m Mode) Validate() error {
	switch m {
	case ModeIsFullOHLC, ModeIsCloseOHLC:
		return nil
	default:
		return ErrInvalidMode
	}
}

// String will return the string representation of the mode.
func (m Mode) String() string {
	return string(m)
}

// Opt will return a pointer to the mode.
func (m Mode) Opt() *Mode {
	return &m
}
