package config

import (
	"errors"
	"fmt"

	"github.com/joho/godotenv"
)

var (
	// ErrInvalidMongoConfig is returned when the MongoDB configuration is invalid.
	ErrInvalidMongoConfig = errors.New("invalid mongo config")
)

// Mongo represents the configuration for the MongoDB database.
type Mongo struct {
	ConnectionString string
	Database         string
}

// LoadMongo loads the MongoDB configuration from the environment variables.
func LoadMongo(defaultValues *Mongo) (c Mongo) {
	if defaultValues != nil {
		c = *defaultValues
	}

	c.setDefault()
	c.overrideFromEnv()

	return c
}

func (c *Mongo) setDefault() {
	overrideString(&c.ConnectionString, "mongodb://localhost:27017")
}

func (c *Mongo) overrideFromEnv() {
	// Attempting to load from .env
	_ = godotenv.Load(".env")

	// Overriding variables
	overrideFromEnv(&c.ConnectionString, "MONGO_CONNECTION_STRING")
	overrideFromEnv(&c.Database, "MONGO_DATABASE")
}

// Validate checks if the configuration is valid.
func (c Mongo) Validate() error {
	if c.ConnectionString == "" {
		return fmt.Errorf(
			"reading connection string from env (%q): %w",
			c.ConnectionString, ErrInvalidMongoConfig)
	}

	if c.Database == "" {
		return fmt.Errorf(
			"reading database from env (%q): %w",
			c.Database, ErrInvalidMongoConfig)
	}

	return nil
}
