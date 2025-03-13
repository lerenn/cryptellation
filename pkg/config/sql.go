package config

import (
	"errors"
	"fmt"

	"github.com/joho/godotenv"
)

// EnvSQLDSN is the environment variable that contains the PostGres DSN.
const EnvSQLDSN = "SQL_DSN"

var (
	// ErrInvalidPostGresConfig is returned when the SQL configuration is invalid.
	ErrInvalidPostGresConfig = errors.New("invalid postgres config")
)

// SQL represents the configuration for the SQL database.
type SQL struct {
	DSN string
}

// LoadSQL loads the SQL configuration from the environment variables.
func LoadSQL(defaultValues *SQL) (c SQL) {
	if defaultValues != nil {
		c = *defaultValues
	}

	c.setDefault()
	c.overrideFromEnv()

	return c
}

func (c *SQL) setDefault() {
	overrideString(&c.DSN, "user=cryptellation password=cryptellation dbname=cryptellation sslmode=disable")
}

func (c *SQL) overrideFromEnv() {
	// Attempting to load from .env
	_ = godotenv.Load(".env")

	// Overriding variables
	overrideFromEnv(&c.DSN, EnvSQLDSN)
}

// Validate checks if the configuration is valid.
func (c SQL) Validate() error {
	if c.DSN == "" {
		return fmt.Errorf("reading DSN from env (%q): %w", c.DSN, ErrInvalidPostGresConfig)
	}

	return nil
}
