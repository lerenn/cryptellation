package daemon

import (
	"context"

	sql "github.com/lerenn/cryptellation/internal/adapters/db/sql/ticks"
	natsTicks "github.com/lerenn/cryptellation/internal/adapters/events/nats/ticks"
	"github.com/lerenn/cryptellation/internal/adapters/exchanges"
	"github.com/lerenn/cryptellation/internal/components/ticks/ports/db"
	"github.com/lerenn/cryptellation/internal/components/ticks/ports/events"
	exchangesPort "github.com/lerenn/cryptellation/internal/components/ticks/ports/exchanges"
	"github.com/lerenn/cryptellation/pkg/config"
)

type adapters struct {
	db        db.Port
	events    events.Port
	exchanges exchangesPort.Port
}

func newAdapters(ctx context.Context) (adapters, error) {
	// Init database client
	db, err := sql.New(config.LoadSQL())
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
