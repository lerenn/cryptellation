package nats

import (
	"context"

	asyncapi "github.com/digital-feather/cryptellation/api/asyncapi/ticks"
	client "github.com/digital-feather/cryptellation/clients/go"
	natsClient "github.com/digital-feather/cryptellation/clients/go/nats"
	"github.com/digital-feather/cryptellation/pkg/config"
	"github.com/digital-feather/cryptellation/pkg/tick"
	"github.com/nats-io/nats.go"
)

type Adapter struct {
	nc     *nats.Conn
	app    *asyncapi.AppController
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
	app, err := asyncapi.NewAppController(asyncapi.NewNATSController(nc))
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

func (a *Adapter) Publish(tick tick.Tick) error {
	// Generated message
	msg := asyncapi.NewTickMessage()
	msg.Payload.Exchange = asyncapi.ExchangeNameSchema(tick.Exchange)
	msg.Payload.PairSymbol = asyncapi.PairSymbolSchema(tick.PairSymbol)
	msg.Payload.Price = tick.Price
	msg.Payload.Time = asyncapi.DateSchema(tick.Time)

	// Send message
	return a.app.PublishCryptellationTicksListenExchangePair(asyncapi.CryptellationTicksListenExchangePairParameters{
		Exchange: asyncapi.ExchangeNameSchema(tick.Exchange),
		Pair:     asyncapi.PairSymbolSchema(tick.PairSymbol),
	}, msg)
}

func (a *Adapter) Subscribe(symbol string) (<-chan tick.Tick, error) {
	return a.client.Listen(context.Background(), client.TicksFilterPayload{
		ExchangeName: "*",
		PairSymbol:   symbol,
	})
}

func (a *Adapter) Close() {
	if a.app != nil {
		a.app.Close()
	}
}
