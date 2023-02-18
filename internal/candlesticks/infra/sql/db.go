package sql

import (
	"fmt"

	"github.com/digital-feather/cryptellation/internal/candlesticks/infra/sql/internal/entities"
	"github.com/digital-feather/cryptellation/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DefaultGormConfig = &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}
)

type DB struct {
	client *gorm.DB
}

func New(c config.SQL) (*DB, error) {
	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("loading sqldb config: %w", err)
	}

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
		&entities.Candlestick{},
	}

	for _, entity := range entities {
		if err := d.client.Debug().Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(entity).Error; err != nil {
			return fmt.Errorf("emptying %T table: %w", entity, err)
		}
	}

	return nil
}
