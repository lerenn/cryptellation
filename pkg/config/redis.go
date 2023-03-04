package config

import (
	"errors"
	"fmt"
	"os"
)

var (
	ErrInvalidRedisConfig = errors.New("invalid redis config")
)

type Redis struct {
	Address  string
	Password string
}

func LoadRedisConfigFromEnv() (r Redis) {
	r.Address = os.Getenv("REDIS_URL")
	r.Password = os.Getenv("REDIS_PASSWORD")

	return r
}

func (c Redis) Validate() error {
	if c.Address == "" {
		return fmt.Errorf("reading address from env (%q): %w", c.Address, ErrInvalidRedisConfig)
	}

	return nil
}
