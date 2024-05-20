package daemon

import (
	"context"

	"github.com/lerenn/cryptellation/pkg/config"
	natsTicks "github.com/lerenn/cryptellation/svc/ticks/internal/adapters/events/nats"
	"github.com/lerenn/cryptellation/svc/ticks/internal/adapters/exchanges"
	"github.com/lerenn/cryptellation/svc/ticks/internal/app/ports/events"
	exchangesPort "github.com/lerenn/cryptellation/svc/ticks/internal/app/ports/exchanges"
)

type adapters struct {
	events    events.Port
	exchanges exchangesPort.Port
}

func newAdapters() (adapters, error) {
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
		events:    events,
		exchanges: exchanges,
	}, nil
}

func (a adapters) Close(ctx context.Context) {
	a.events.Close(ctx)
}
