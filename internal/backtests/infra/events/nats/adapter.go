//go:generate asyncapi-codegen -g application -p generated -i ../../../../../api/asyncapi-spec/backtests.yaml -o ./generated/app.gen.go
//go:generate asyncapi-codegen -g client      -p generated -i ../../../../../api/asyncapi-spec/backtests.yaml -o ./generated/client.gen.go
//go:generate asyncapi-codegen -g broker      -p generated -i ../../../../../api/asyncapi-spec/backtests.yaml -o ./generated/broker.gen.go
//go:generate asyncapi-codegen -g types       -p generated -i ../../../../../api/asyncapi-spec/backtests.yaml -o ./generated/types.gen.go
//go:generate asyncapi-codegen -g nats        -p generated -i ../../../../../api/asyncapi-spec/backtests.yaml -o ./generated/nats.gen.go

package nats

import (
	"context"

	client "github.com/digital-feather/cryptellation/clients/go"
	natsClient "github.com/digital-feather/cryptellation/clients/go/nats"
	"github.com/digital-feather/cryptellation/internal/backtests/infra/events/nats/generated"
	"github.com/digital-feather/cryptellation/pkg/config"
	"github.com/digital-feather/cryptellation/pkg/types/event"
	"github.com/digital-feather/cryptellation/pkg/types/tick"
	"github.com/nats-io/nats.go"
)

type Adapter struct {
	nc     *nats.Conn
	app    *generated.AppController
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
	app, err := generated.NewAppController(generated.NewNATSController(nc))
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
	msg := generated.NewBacktestsEventMessage()
	msg.Payload.Time = generated.DateSchema(evt.Time)
	msg.Payload.Type = evt.Type.String()

	// Set message depending on event type
	switch evt.Type {
	case event.TypeIsStatus:
		statusEvt, ok := evt.Content.(event.Status)
		if !ok {
			return event.ErrMismatchingType
		}

		msg.Payload.Content.Finished = statusEvt.Finished
	case event.TypeIsTick:
		t, ok := evt.Content.(tick.Tick)
		if !ok {
			return event.ErrMismatchingType
		}

		msg.Payload.Content.Exchange = generated.ExchangeNameSchema(t.Exchange)
		msg.Payload.Content.PairSymbol = generated.PairSymbolSchema(t.PairSymbol)
		msg.Payload.Content.Price = t.Price
		msg.Payload.Content.Time = generated.DateSchema(t.Time)
	default:
		return event.ErrUnknownType
	}

	// Send message
	return a.app.PublishBacktestsEventsID(generated.BacktestsEventsIDParameters{
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
