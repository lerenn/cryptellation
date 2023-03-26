package sql

import (
	"fmt"
	"log"

	"github.com/digital-feather/cryptellation/pkg/config"
	"github.com/digital-feather/cryptellation/services/exchanges/io/db/adapters/sql/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	client *gorm.DB
}

func New(c config.SQL) (*DB, error) {
	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("loading sqldb config: %w", err)
	}

	client, err := gorm.Open(postgres.Open(c.URL()), config.DefaultGormConfig)
	if err != nil {
		return nil, fmt.Errorf("opening sqldb connection: %w", err)
	}

	db := &DB{
		client: client,
	}

	return db, nil
}

func (sqldb *DB) Reset() error {
	entities := []interface{}{
		&entities.Pair{},
		&entities.Period{},
		&entities.Exchange{},
	}

	sessionOpt := &gorm.Session{AllowGlobalUpdate: true}
	for i, entity := range entities {
		log.Println(i)
		if err := sqldb.client.Session(sessionOpt).Delete(entity).Error; err != nil {
			return fmt.Errorf("emptying %T table: %w", entity, err)
		}
	}

	return nil
}
