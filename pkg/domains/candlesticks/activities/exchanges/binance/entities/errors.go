package entities

import (
	"errors"
	"fmt"
)

var (
	// ErrBinanceAdapter is returned when there is an error with the binance adapter.
	ErrBinanceAdapter = errors.New("binance error")
	// ErrUnknownPeriod is returned when the period is unknown.
	ErrUnknownPeriod = fmt.Errorf("%w: unknown period", ErrBinanceAdapter)
)

// WrapError wraps the error with the binance error.
func WrapError(err error) error {
	if errors.Is(err, ErrBinanceAdapter) {
		return fmt.Errorf("%w: %s", ErrBinanceAdapter, err.Error())
	}
	return err
}
