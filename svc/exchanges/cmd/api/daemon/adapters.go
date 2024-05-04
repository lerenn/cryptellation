package daemon

import (
	"context"

	"github.com/lerenn/cryptellation/pkg/config"
	sql "github.com/lerenn/cryptellation/svc/exchanges/internal/adapters/db/sql"
	"github.com/lerenn/cryptellation/svc/exchanges/internal/adapters/exchanges"
	"github.com/lerenn/cryptellation/svc/exchanges/internal/app/ports/db"
	exchangesPort "github.com/lerenn/cryptellation/svc/exchanges/internal/app/ports/exchanges"
)

type adapters struct {
	db        db.Port
	exchanges exchangesPort.Port
}

func newAdapters(ctx context.Context) (adapters, error) {
	// Init database client
	db, err := sql.New(config.LoadSQL(nil))
	if err != nil {
		return adapters{}, err
	}

	// Init exchanges connections
	exchanges, err := exchanges.New()
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
