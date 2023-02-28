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

func LoadNATSConfigFromEnv() (n NATS) {
	n.OverrideFromEnv()
	return n
}

func LoadDefaultNATSConfig() (n NATS) {
	n.LoadDefault()
	return n
}

func (n NATS) URL() string {
	return fmt.Sprintf("nats://%s:%d", n.Host, n.Port)
}

func (n *NATS) OverrideFromEnv() {
	host := os.Getenv("NATS_HOST")
	if host != "" {
		n.Host = host
	}

	port, _ := strconv.Atoi(os.Getenv("NATS_PORT"))
	if port != 0 {
		n.Port = port
	}
}

func (n *NATS) LoadDefault() {
	n.Host = "localhost"
	n.Port = 4222
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
