package config

import (
	"errors"
	"fmt"

	"github.com/joho/godotenv"
)

const EnvPostgresDSN = "POSTGRES_DSN"

var (
	// ErrInvalidPostGresConfig is returned when the PostGresDB configuration is invalid.
	ErrInvalidPostGresConfig = errors.New("invalid postgres config")
)

// PostGres represents the configuration for the PostGresDB database.
type PostGres struct {
	DSN string
}

// LoadPostGres loads the PostGresDB configuration from the environment variables.
func LoadPostGres(defaultValues *PostGres) (c PostGres) {
	if defaultValues != nil {
		c = *defaultValues
	}

	c.setDefault()
	c.overrideFromEnv()

	return c
}

func (c *PostGres) setDefault() {
	overrideString(&c.DSN, "user=cryptellation password=cryptellation dbname=cryptellation sslmode=disable")
}

func (c *PostGres) overrideFromEnv() {
	// Attempting to load from .env
	_ = godotenv.Load(".env")

	// Overriding variables
	overrideFromEnv(&c.DSN, EnvPostgresDSN)
}

// Validate checks if the configuration is valid.
func (c PostGres) Validate() error {
	if c.DSN == "" {
		return fmt.Errorf("reading DSN from env (%q): %w", c.DSN, ErrInvalidPostGresConfig)
	}

	return nil
}
