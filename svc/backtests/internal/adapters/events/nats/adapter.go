package backtests

import (
	adapter "cryptellation/internal/adapters/events/nats"
	pkg "cryptellation/internal/asyncapi"
	"cryptellation/pkg/config"

	asyncapi "cryptellation/svc/backtests/api/asyncapi"
	client "cryptellation/svc/backtests/clients/go"
	backtestsnats "cryptellation/svc/backtests/clients/go/nats"
)

type Adapter struct {
	events *adapter.Adapter
	app    *asyncapi.AppController
	client client.Client
}

func New(c config.NATS) (*Adapter, error) {
	// Create embedded database access
	events, err := adapter.New(c)
	if err != nil {
		return nil, err
	}

	// Create new app controller
	app, err := asyncapi.NewAppController(events.Broker(), asyncapi.WithLogger(pkg.LoggerWrapper{}))
	if err != nil {
		return nil, err
	}

	// Create a new client
	client, err := backtestsnats.New(c, backtestsnats.WithLogger(pkg.LoggerWrapper{}))
	if err != nil {
		return nil, err
	}

	return &Adapter{
		events: events,
		app:    app,
		client: client,
	}, nil
}
