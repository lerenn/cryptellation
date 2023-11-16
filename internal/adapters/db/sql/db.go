package sql

import (
	"fmt"

	"github.com/lerenn/cryptellation/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Adapter struct {
	Client *gorm.DB
}

func New(c config.SQL) (*Adapter, error) {
	// Validate configuration
	if err := c.Validate(); err != nil {
		return nil, err
	}

	// Generate SQL client
	client, err := gorm.Open(postgres.Open(c.URL()), config.DefaultGormConfig)
	if err != nil {
		return nil, fmt.Errorf("opening sqldb connection: %w", err)
	}

	// Return client
	return &Adapter{
		Client: client,
	}, nil
}
