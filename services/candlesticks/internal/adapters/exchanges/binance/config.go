package binance

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidConfig = errors.New("invalid binance config")
)

type Config struct {
	ApiKey    string
	SecretKey string
}

func (c Config) Validate() error {
	if c.ApiKey == "" {
		return fmt.Errorf("reading api key from env (%q): %w", c.ApiKey, ErrInvalidConfig)
	}

	if c.SecretKey == "" {
		return fmt.Errorf("reading secret key from env (%q): %w", c.SecretKey, ErrInvalidConfig)
	}

	return nil
}
