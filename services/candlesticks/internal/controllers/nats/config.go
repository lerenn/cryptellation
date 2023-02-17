package nats

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidConfig = errors.New("invalid nats config")
)

type Config struct {
	Host string
	Port int
}

func (c Config) URL() string {
	return fmt.Sprintf("nats://%s:%d", c.Host, c.Port)
}

func (c Config) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("reading host from env (%q): %w", c.Host, ErrInvalidConfig)
	}

	if c.Port == 0 {
		return fmt.Errorf("reading port from env (%q): %w", c.Port, ErrInvalidConfig)
	}

	return nil
}
