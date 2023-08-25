package nats

import (
	"context"

	client "github.com/lerenn/cryptellation/clients/go"
	natsClient "github.com/lerenn/cryptellation/clients/go/nats"
	"github.com/lerenn/cryptellation/internal/ctrl/ticks/events"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/models/tick"
	"github.com/nats-io/nats.go"
)

type Adapter struct {
	nc     *nats.Conn
	app    *events.AppController
	client client.Ticks
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
	app, err := events.NewAppController(events.NewNATSController(nc))
	if err != nil {
		return nil, err
	}

	// Create a new client
	client, err := natsClient.NewTicks(c)
	if err != nil {
		return nil, err
	}

	return &Adapter{
		nc:     nc,
		app:    app,
		client: client,
	}, nil
}

func (a *Adapter) Publish(ctx context.Context, tick tick.Tick) error {
	// Generated message
	msg := events.NewTickMessage()
	msg.Payload.Exchange = events.ExchangeNameSchema(tick.Exchange)
	msg.Payload.PairSymbol = events.PairSymbolSchema(tick.PairSymbol)
	msg.Payload.Price = tick.Price
	msg.Payload.Time = events.DateSchema(tick.Time)

	// Send message
	return a.app.PublishCryptellationTicksListenExchangePair(ctx,
		events.CryptellationTicksListenExchangePairParameters{
			Exchange: events.ExchangeNameSchema(tick.Exchange),
			Pair:     events.PairSymbolSchema(tick.PairSymbol),
		}, msg)
}

func (a *Adapter) Subscribe(ctx context.Context, symbol string) (<-chan tick.Tick, error) {
	return a.client.Listen(ctx, client.TicksFilterPayload{
		ExchangeName: "*",
		PairSymbol:   symbol,
	})
}

func (a *Adapter) Close(ctx context.Context) {
	if a.app != nil {
		a.app.Close(ctx)
	}
}
