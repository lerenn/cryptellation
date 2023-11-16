package config

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidNATS = errors.New("invalid nats config")
)

type NATS struct {
	Host string
	Port int
}

func LoadNATS() (c NATS) {
	c.setDefault()
	c.overrideFromEnv()
	return c
}

func (c *NATS) setDefault() {
	c.Host = "localhost"
	c.Port = 4222
}

func (c *NATS) overrideFromEnv() {
	overrideFromEnv(&c.Host, "NATS_HOST")
	overrideIntFromEnv(&c.Port, "NATS_PORT")
}

func (c NATS) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("reading host from env (%q): %w", c.Host, ErrInvalidNATS)
	}

	if c.Port == 0 {
		return fmt.Errorf("reading port from env (%q): %w", c.Port, ErrInvalidNATS)
	}

	return nil
}

func (c NATS) URL() string {
	return fmt.Sprintf("nats://%s:%d", c.Host, c.Port)
}
