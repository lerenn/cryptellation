package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	// ErrInvalidBinance is returned when the Binance configuration is invalid.
	ErrInvalidBinance = errors.New("invalid binance config")
)

// Binance represents the configuration for the Binance exchange.
type Binance struct {
	APIKey    string
	SecretKey string
}

// LoadBinance loads the Binance configuration from the environment variables.
func LoadBinance(defaultValues *Binance, additionalEnvFilePaths ...string) (c Binance) {
	if defaultValues != nil {
		c = *defaultValues
	}

	c.setDefault()
	c.overrideFromEnv(additionalEnvFilePaths...)

	return c
}

// LoadBinanceTest loads the Binance configuration from the environment variables for testing purposes.
func LoadBinanceTest() Binance {
	// Find the .credentials.env file in the parent directories
	p := ".credentials.env"
	for i := 0; i < 20; i++ {
		if _, err := os.Stat(p); err == nil {
			break
		}

		p = "../" + p
	}

	// Load the config
	return LoadBinance(nil, p)
}

func (c *Binance) setDefault() {
	// Nothing to do
}

func (c *Binance) overrideFromEnv(additionalEnvFilePaths ...string) {
	// Attempting to load from .env
	_ = godotenv.Load(".env")
	_ = godotenv.Load(additionalEnvFilePaths...)

	// Overriding variables
	overrideFromEnv(&c.APIKey, "BINANCE_API_KEY")
	overrideFromEnv(&c.SecretKey, "BINANCE_SECRET_KEY")
}

// Validate checks if the configuration is valid.
func (c Binance) Validate() error {
	if c.APIKey == "" {
		return fmt.Errorf("reading api key from env (%q): %w", c.APIKey, ErrInvalidBinance)
	}

	if c.SecretKey == "" {
		return fmt.Errorf("reading secret key from env (%q): %w", c.SecretKey, ErrInvalidBinance)
	}

	return nil
}
