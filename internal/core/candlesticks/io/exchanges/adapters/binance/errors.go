package binance

import (
	"errors"
	"fmt"
)

var (
	ErrBinanceAdapter = errors.New("binance adapter error")

	ErrUnknownPeriod = fmt.Errorf("%w: unknown period", ErrBinanceAdapter)
)

func wrapError(err error) error {
	if errors.Is(err, ErrBinanceAdapter) {
		return fmt.Errorf("%w: %s", ErrBinanceAdapter, err.Error())
	}
	return err
}
