package nats

import (
	"errors"
	"fmt"
	"os"
)

var (
	ErrInvalidConfig = errors.New("invalid NATS config")
)

type Config struct {
	URL string
}

func (c *Config) Load() *Config {
	c.URL = os.Getenv("NATS_URL")
	return c
}

func (c Config) Validate() error {
	if c.URL == "" {
		return fmt.Errorf("reading URL from env (%q): %w", c.URL, ErrInvalidConfig)
	}

	return nil
}
