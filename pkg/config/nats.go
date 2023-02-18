package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

var (
	ErrInvalidNATS = errors.New("invalid nats config")
)

type NATS struct {
	Host string
	Port int
}

func LoadNATSConfigFromEnv() (c NATS) {
	c.Host = os.Getenv("NATS_HOST")
	c.Port, _ = strconv.Atoi(os.Getenv("NATS_PORT"))
	return c
}

func (n NATS) URL() string {
	return fmt.Sprintf("nats://%s:%d", n.Host, n.Port)
}

func (n NATS) Validate() error {
	if n.Host == "" {
		return fmt.Errorf("reading host from env (%q): %w", n.Host, ErrInvalidNATS)
	}

	if n.Port == 0 {
		return fmt.Errorf("reading port from env (%q): %w", n.Port, ErrInvalidNATS)
	}

	return nil
}
