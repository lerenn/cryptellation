package sql

import (
	"fmt"

	"github.com/lerenn/cryptellation/internal/infra/backtests/db/sql/entities"
	"github.com/lerenn/cryptellation/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	client *gorm.DB
}

func New(c config.SQL) (*DB, error) {
	// Validate configuration
	if err := c.Validate(); err != nil {
		return nil, err
	}

	// Generate SQL client
	client, err := gorm.Open(postgres.Open(c.URL()), DefaultGormConfig)
	if err != nil {
		return nil, fmt.Errorf("opening sqldb connection: %w", err)
	}

	// Return client
	return &DB{
		client: client,
	}, nil
}

func (d *DB) Reset() error {
	entities := []interface{}{
		&entities.Balance{},
		&entities.Backtest{},
		&entities.Order{},
		&entities.TickSubscription{},
	}

	for _, entity := range entities {
		if err := d.client.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(entity).Error; err != nil {
			return fmt.Errorf("emptying %T table: %w", entity, err)
		}
	}

	return nil
}
