package sql

import (
	"fmt"

	"github.com/lerenn/cryptellation/internal/infra/ticks/db/sql/entities"
	"github.com/lerenn/cryptellation/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	client *gorm.DB
}

func New(cfg config.SQL) (*DB, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("loading sqldb config: %w", err)
	}

	client, err := gorm.Open(postgres.Open(cfg.URL()), config.DefaultGormConfig)
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
		&entities.SymbolListener{},
	}

	for _, entity := range entities {
		if err := d.client.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(entity).Error; err != nil {
			return fmt.Errorf("emptying %T table: %w", entity, err)
		}
	}

	return nil
}
