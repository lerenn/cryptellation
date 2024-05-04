package daemon

import (
	"context"

	"github.com/lerenn/cryptellation/pkg/config"
	sql "github.com/lerenn/cryptellation/svc/ticks/internal/adapters/db/sql"
	natsTicks "github.com/lerenn/cryptellation/svc/ticks/internal/adapters/events/nats"
	"github.com/lerenn/cryptellation/svc/ticks/internal/adapters/exchanges"
	"github.com/lerenn/cryptellation/svc/ticks/internal/app/ports/db"
	"github.com/lerenn/cryptellation/svc/ticks/internal/app/ports/events"
	exchangesPort "github.com/lerenn/cryptellation/svc/ticks/internal/app/ports/exchanges"
)

type adapters struct {
	db        db.Port
	events    events.Port
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

	// Init Events client
	events, err := natsTicks.New(config.LoadNATS())
	if err != nil {
		return adapters{}, err
	}

	return adapters{
		db:        db,
		events:    events,
		exchanges: exchanges,
	}, nil
}

func (a adapters) Close(ctx context.Context) {
	a.events.Close(ctx)
}
