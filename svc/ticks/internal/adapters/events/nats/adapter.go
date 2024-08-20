package nats

import (
	adapter "cryptellation/internal/adapters/events/nats"
	"cryptellation/pkg/config"

	asyncapi "cryptellation/svc/ticks/api/asyncapi"
)

type Adapter struct {
	events *adapter.Adapter
	app    *asyncapi.AppController
	user   *asyncapi.UserController
}

func New(c config.NATS) (*Adapter, error) {
	// Create embedded database access
	events, err := adapter.New(c)
	if err != nil {
		return nil, err
	}

	// Create new app controller
	app, err := asyncapi.NewAppController(events.Broker())
	if err != nil {
		return nil, err
	}

	// Create a new user controller
	user, err := asyncapi.NewUserController(events.Broker())
	if err != nil {
		return nil, err
	}

	return &Adapter{
		events: events,
		app:    app,
		user:   user,
	}, nil
}
