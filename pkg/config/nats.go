package config

import (
	"errors"
	"fmt"

	"github.com/joho/godotenv"
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

const (
	NatsHostEnvName = "NATS_HOST"
	NatsPortEnvName = "NATS_PORT"
)

func (c *NATS) overrideFromEnv() {
	// Attempting to load from .env
	_ = godotenv.Load(".env")

	// Overriding variables
	overrideFromEnv(&c.Host, NatsHostEnvName)
	overrideIntFromEnv(&c.Port, NatsPortEnvName)
}

func (c NATS) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("NATS host is empty (%q env is %q): %w", NatsHostEnvName, c.Host, ErrInvalidNATS)
	}

	if c.Port == 0 {
		return fmt.Errorf("NATS port is empty (%q env is %q): %w", NatsPortEnvName, c.Port, ErrInvalidNATS)
	}

	return nil
}

func (c NATS) URL() string {
	return fmt.Sprintf("nats://%s:%d", c.Host, c.Port)
}
