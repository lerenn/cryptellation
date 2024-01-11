package entities

import (
	"errors"
	"fmt"
)

var (
	ErrBinanceAdapter = errors.New("binance error")
	ErrUnknownPeriod  = fmt.Errorf("%w: unknown period", ErrBinanceAdapter)
)

func WrapError(err error) error {
	if errors.Is(err, ErrBinanceAdapter) {
		return fmt.Errorf("%w: %s", ErrBinanceAdapter, err.Error())
	}
	return err
}
