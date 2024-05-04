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

func LoadRedis(defaultValues *Redis) (c Redis) {
	if defaultValues != nil {
		c = *defaultValues
	}

	c.setDefault()
	c.overrideFromEnv()

	return c
}

func (c *Redis) setDefault() {
	overrideString(&c.Address, "localhost:6379")
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
