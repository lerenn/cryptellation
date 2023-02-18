package config

import (
	"errors"
	"fmt"
	"os"
)

var (
	ErrInvalidBinance = errors.New("invalid binance config")
)

type Binance struct {
	ApiKey    string
	SecretKey string
}

func LoadBinanceConfigFromEnv() (c Binance) {
	c.ApiKey = os.Getenv("BINANCE_API_KEY")
	c.SecretKey = os.Getenv("BINANCE_SECRET_KEY")
	return c
}

func (b Binance) Validate() error {
	if b.ApiKey == "" {
		return fmt.Errorf("reading api key from env (%q): %w", b.ApiKey, ErrInvalidBinance)
	}

	if b.SecretKey == "" {
		return fmt.Errorf("reading secret key from env (%q): %w", b.SecretKey, ErrInvalidBinance)
	}

	return nil
}
