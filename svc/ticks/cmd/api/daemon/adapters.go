package daemon

import (
	"context"

	"cryptellation/internal/config"

	natsTicks "cryptellation/svc/ticks/internal/adapters/events/nats"
	"cryptellation/svc/ticks/internal/adapters/exchanges"
	"cryptellation/svc/ticks/internal/app/ports/events"
	exchangesPort "cryptellation/svc/ticks/internal/app/ports/exchanges"
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
