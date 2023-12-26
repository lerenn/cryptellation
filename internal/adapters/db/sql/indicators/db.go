package indicators

import (
	"context"

	"github.com/lerenn/cryptellation/internal/adapters/db/sql"
	adapter "github.com/lerenn/cryptellation/internal/adapters/db/sql"
	"github.com/lerenn/cryptellation/internal/adapters/db/sql/indicators/entities"
	"github.com/lerenn/cryptellation/pkg/config"
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

func (a *Adapter) Reset(ctx context.Context) error {
	return sql.Reset(ctx, a.db.Client, []interface{}{
		&entities.SimpleMovingAverage{},
	})
}
