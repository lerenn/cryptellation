package sql

import (
	"fmt"

	"github.com/digital-feather/cryptellation/internal/backtests/infra/db/sql/entities"
	"github.com/digital-feather/cryptellation/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	client *gorm.DB
}

func New(c config.SQL) (*DB, error) {
	client, err := gorm.Open(postgres.Open(c.URL()), DefaultGormConfig)
	if err != nil {
		return nil, fmt.Errorf("opening sqldb connection: %w", err)
	}

	db := &DB{
		client: client,
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
