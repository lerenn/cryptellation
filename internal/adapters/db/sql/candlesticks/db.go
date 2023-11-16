package candlesticks

import (
	"fmt"

	adapter "github.com/lerenn/cryptellation/internal/adapters/db/sql"
	"github.com/lerenn/cryptellation/internal/adapters/db/sql/candlesticks/entities"
	"github.com/lerenn/cryptellation/pkg/config"
	"gorm.io/gorm"
)

type Adapter struct {
	db *adapter.Adapter
}

func New(c config.SQL) (*Adapter, error) {
	// Create embedded database access
	db, err := adapter.New(c)

	// Return database access
	return &Adapter{
		db: db,
	}, err
}

func (a *Adapter) Reset() error {
	entities := []interface{}{
		&entities.Candlestick{},
	}

	for _, entity := range entities {
		if err := a.db.Client.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(entity).Error; err != nil {
			return fmt.Errorf("emptying %T table: %w", entity, err)
		}
	}

	return nil
}
