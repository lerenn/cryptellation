package config

import (
	"errors"
	"fmt"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	ErrInvalidSQL = errors.New("invalid sql config")
)

var (
	DefaultGormConfig = &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}
)

type SQL struct {
	User     string
	Password string
	Host     string
	Port     int
	Database string
}

func LoadSQL() (c SQL) {
	c.setDefault()
	c.overrideFromEnv()
	return c
}

func (c *SQL) setDefault() {
	// Nothing to do
}

func (c *SQL) overrideFromEnv() {
	// Attempting to load from .env
	_ = godotenv.Load(".env")

	// Overriding variables
	overrideFromEnv(&c.Host, "SQLDB_HOST")
	overrideIntFromEnv(&c.Port, "SQLDB_PORT")
	overrideFromEnv(&c.User, "SQLDB_USER")
	overrideFromEnv(&c.Password, "SQLDB_PASSWORD")
	overrideFromEnv(&c.Database, "SQLDB_DATABASE")
}

func (c SQL) URL() string {
	fields := map[string]string{
		"host":     c.Host,
		"port":     fmt.Sprintf("%d", c.Port),
		"user":     c.User,
		"password": c.Password,
		"dbname":   c.Database,
	}

	var dsn string
	for field, value := range fields {
		if value != "" {
			dsn = fmt.Sprintf("%s %s=%s", dsn, field, value)
		}
	}

	return dsn
}

func (c SQL) Validate() error {
	if c.User == "" {
		return fmt.Errorf("reading user from env (%q): %w", c.User, ErrInvalidSQL)
	}

	if c.Host == "" {
		return fmt.Errorf("reading host from env (%q): %w", c.Host, ErrInvalidSQL)
	}

	if c.Port == 0 {
		return fmt.Errorf("reading port from env (%q): %w", c.Port, ErrInvalidSQL)
	}

	if c.Database == "" {
		return fmt.Errorf("reading database from env (%q): %w", c.Database, ErrInvalidSQL)
	}

	return nil
}
