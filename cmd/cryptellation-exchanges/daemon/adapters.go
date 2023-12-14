package daemon

import (
	"context"

	sql "github.com/lerenn/cryptellation/internal/adapters/db/sql/exchanges"
	exchangesAdapter "github.com/lerenn/cryptellation/internal/adapters/exchanges"
	"github.com/lerenn/cryptellation/internal/components/exchanges/ports/db"
	exchangesPort "github.com/lerenn/cryptellation/internal/components/exchanges/ports/exchanges"
	"github.com/lerenn/cryptellation/pkg/config"
)

type adapters struct {
	db        db.Port
	exchanges exchangesPort.Port
}

func newAdapters(ctx context.Context) (adapters, error) {
	// Init database client
	db, err := sql.New(config.LoadSQL())
	if err != nil {
		return adapters{}, err
	}

	// Init exchanges connections
	exchanges, err := exchangesAdapter.New()
	if err != nil {
		return adapters{}, err
	}

	return adapters{
		db:        db,
		exchanges: exchanges,
	}, nil
}

func (a adapters) Close(ctx context.Context) {
}
