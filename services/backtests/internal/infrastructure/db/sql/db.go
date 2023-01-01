package sql

import (
	"fmt"

	"github.com/digital-feather/cryptellation/services/backtests/internal/infrastructure/db/sql/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	client *gorm.DB
	config Config
}

func New() (*DB, error) {
	var c Config
	if err := c.Load().Validate(); err != nil {
		return nil, fmt.Errorf("loading sqldb config: %w", err)
	}

	client, err := gorm.Open(postgres.Open(c.URL()), DefaultGormConfig)
	if err != nil {
		return nil, fmt.Errorf("opening sqldb connection: %w", err)
	}

	db := &DB{
		client: client,
		config: c,
	}

	return db, nil
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
