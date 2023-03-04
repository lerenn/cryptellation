package nats

import (
	"context"

	serviceNATSClient "github.com/digital-feather/cryptellation/internal/ticks/ctrl/nats"
	"github.com/digital-feather/cryptellation/internal/ticks/ctrl/nats/generated"
	"github.com/digital-feather/cryptellation/pkg/config"
	"github.com/digital-feather/cryptellation/pkg/tick"
	"github.com/nats-io/nats.go"
)

type Adapter struct {
	nc     *nats.Conn
	app    *generated.AppController
	client serviceNATSClient.Client
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
	app, err := generated.NewAppController(generated.NewNATSController(nc))
	if err != nil {
		return nil, err
	}

	// Create a new client
	client, err := serviceNATSClient.New(c)
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
	msg := generated.NewTickMessage()
	msg.Payload.Exchange = generated.ExchangeNameSchema(tick.Exchange)
	msg.Payload.PairSymbol = generated.PairSymbolSchema(tick.PairSymbol)
	msg.Payload.Price = tick.Price
	msg.Payload.Time = generated.DateSchema(tick.Time)

	// Send message
	return a.app.PublishTicksListenExchangePair(generated.TicksListenExchangePairParameters{
		Exchange: tick.Exchange,
		Pair:     tick.PairSymbol,
	}, msg)
}

func (a *Adapter) Subscribe(symbol string) (<-chan tick.Tick, error) {
	return a.client.Listen(context.Background(), serviceNATSClient.TicksFilterPayload{
		ExchangeName: "*",
		PairSymbol:   symbol,
	})
}

func (a *Adapter) Close() {
	if a.app != nil {
		a.app.Close()
	}
}
