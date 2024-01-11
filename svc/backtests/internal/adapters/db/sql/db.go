package backtests

import (
	"context"

	"github.com/lerenn/cryptellation/pkg/adapters/db/sql"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/svc/backtests/internal/adapters/db/sql/entities"
)

type Adapter struct {
	db *sql.Adapter
}

func New(c config.SQL) (*Adapter, error) {
	// Create embedded database access
	db, err := sql.New(c)

	// Return database access
	return &Adapter{
		db: db,
	}, err
}

func (a *Adapter) Reset(ctx context.Context) error {
	return sql.Reset(ctx, a.db.Client, []interface{}{
		&entities.Balance{},
		&entities.Backtest{},
		&entities.Order{},
		&entities.TickSubscription{},
	})
}
