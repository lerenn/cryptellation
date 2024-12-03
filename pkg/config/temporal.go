package config

import (
	"errors"
	"fmt"

	"github.com/joho/godotenv"
	temporalclient "go.temporal.io/sdk/client"
)

const (
	// TemporalAddressEnvName is the name of the environment variable that
	// contains the address of the temporal server.
	TemporalAddressEnvName = "TEMPORAL_ADDRESS"
)

var (
	// ErrInvalidTemporalConfig is returned when the temporal configuration is invalid.
	ErrInvalidTemporalConfig = errors.New("invalid temporal config")
)

// Temporal contains the configuration for the temporal server.
type Temporal struct {
	Address string
}

// LoadTemporal loads the temporal configuration from the environment.
func LoadTemporal(defaultValues *Temporal) (c Temporal) {
	if defaultValues != nil {
		c = *defaultValues
	}

	c.setDefault()
	c.overrideFromEnv()

	return c
}

// CreateTemporalClient creates a new temporal client from the configuration.
func (c Temporal) CreateTemporalClient() (temporalclient.Client, error) {
	return temporalclient.Dial(temporalclient.Options{
		HostPort: c.Address,
	})
}

func (c *Temporal) setDefault() {
	overrideString(&c.Address, "localhost:7233")
}

func (c *Temporal) overrideFromEnv() {
	// Attempting to load from .env
	_ = godotenv.Load(".env")

	// Overriding variables
	overrideFromEnv(&c.Address, TemporalAddressEnvName)
}

// Validate checks if the configuration is valid.
func (c Temporal) Validate() error {
	if c.Address == "" {
		return fmt.Errorf(
			"%w: reading address from env (%s=%q)",
			ErrInvalidTemporalConfig, TemporalAddressEnvName, c.Address)
	}

	return nil
}
