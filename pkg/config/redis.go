package config

import (
	"errors"
	"fmt"

	"github.com/joho/godotenv"
)

var (
	ErrInvalidConfig = errors.New("invalid redis config")
)

type Redis struct {
	Address  string
	Password string
}

func LoadRedis() (c Redis) {
	c.setDefault()
	c.overrideFromEnv()
	return c
}

func (c *Redis) setDefault() {
	// Nothing to do
}

func (c *Redis) overrideFromEnv() {
	// Attempting to load from .env
	_ = godotenv.Load(".env")

	// Overriding variables
	overrideFromEnv(&c.Address, "REDIS_URL")
	overrideFromEnv(&c.Password, "REDIS_PASSWORD")
}

func (c Redis) Validate() error {
	if c.Address == "" {
		return fmt.Errorf("reading address from env (%q): %w", c.Address, ErrInvalidConfig)
	}

	return nil
}
