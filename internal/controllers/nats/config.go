package nats

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

var (
	ErrInvalidConfig = errors.New("invalid nats config")
)

type Config struct {
	Host string
	Port int
}

func LoadConfigFromEnv() (c Config) {
	c.Host = os.Getenv("NATS_HOST")
	c.Port, _ = strconv.Atoi(os.Getenv("NATS_PORT"))
	return c
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
