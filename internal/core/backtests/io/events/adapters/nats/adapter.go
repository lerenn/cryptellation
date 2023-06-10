package nats

import (
	"context"

	asyncapi "github.com/lerenn/cryptellation/api/asyncapi/backtests"
	client "github.com/lerenn/cryptellation/clients/go"
	natsClient "github.com/lerenn/cryptellation/clients/go/nats"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/event"
	"github.com/nats-io/nats.go"
)

type Adapter struct {
	nc     *nats.Conn
	app    *asyncapi.AppController
	client client.Backtests
}

func New(c config.NATS) (*Adapter, error) {
	// Validate configuration
	if err := c.Validate(); err != nil {
		return nil, err
	}

	// Connect to NATS
	nc, err := nats.Connect(c.URL())
	if err != nil {
		return nil, err
	}

	// Create new app controller
	app, err := asyncapi.NewAppController(asyncapi.NewNATSController(nc))
	if err != nil {
		return nil, err
	}

	// Create a new client
	client, err := natsClient.NewBacktests(c)
	if err != nil {
		return nil, err
	}

	return &Adapter{
		nc:     nc,
		app:    app,
		client: client,
	}, nil
}

func (a *Adapter) Publish(backtestID uint, evt event.Event) error {
	// Generated message
	msg := asyncapi.NewBacktestsEventMessage()

	// Set from event
	if err := msg.Set(evt); err != nil {
		return err
	}

	// Send message
	return a.app.PublishCryptellationBacktestsEventsID(asyncapi.CryptellationBacktestsEventsIDParameters{
		ID: int64(backtestID),
	}, msg)
}

func (a *Adapter) Subscribe(backtestID uint) (<-chan event.Event, error) {
	return a.client.ListenEvents(context.Background(), backtestID)
}

func (a *Adapter) Close() {
	if a.app != nil {
		a.app.Close()
	}
}
