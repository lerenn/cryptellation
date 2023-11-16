package config

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidBinance = errors.New("invalid binance config")
)

type Binance struct {
	ApiKey    string
	SecretKey string
}

func LoadBinance() (c Binance) {
	c.setDefault()
	c.overrideFromEnv()
	return c
}

func (c *Binance) setDefault() {
	// Nothing to do
}

func (c *Binance) overrideFromEnv() {
	overrideFromEnv(&c.ApiKey, "BINANCE_API_KEY")
	overrideFromEnv(&c.SecretKey, "BINANCE_SECRET_KEY")
}

func (c Binance) Validate() error {
	if c.ApiKey == "" {
		return fmt.Errorf("reading api key from env (%q): %w", c.ApiKey, ErrInvalidBinance)
	}

	if c.SecretKey == "" {
		return fmt.Errorf("reading secret key from env (%q): %w", c.SecretKey, ErrInvalidBinance)
	}

	return nil
}
