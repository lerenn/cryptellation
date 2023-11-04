package nats

import (
	"context"
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/loggers"
	client "github.com/lerenn/cryptellation/clients/go"
	asyncapi "github.com/lerenn/cryptellation/pkg/asyncapi/ticks"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/models/tick"
)

type Ticks struct {
	broker *nats.Controller
	ctrl   *asyncapi.UserController
	logger extensions.Logger
}

func NewTicks(c config.NATS) (client.Ticks, error) {
	// Create a NATS Controller
	broker := nats.NewController(c.URL())

	// Create a logger
	logger := loggers.NewECS()

	// Create a new user controller
	ctrl, err := asyncapi.NewUserController(broker, asyncapi.WithLogger(logger))
	if err != nil {
		return nil, err
	}

	return Ticks{
		broker: broker,
		ctrl:   ctrl,
		logger: logger,
	}, nil
}

func (t Ticks) Register(ctx context.Context, payload client.TicksFilterPayload) error {
	// Generate message
	msg := asyncapi.NewRegisteringRequestMessage()
	msg.Set(payload)

	// Send message
	resp, err := t.ctrl.WaitForCryptellationTicksRegisterResponse(ctx, &msg, func(ctx context.Context) error {
		return t.ctrl.PublishCryptellationTicksRegisterRequest(ctx, msg)
	})
	if err != nil {
		return err
	}

	// Check error from server
	if resp.Payload.Error != nil {
		return fmt.Errorf("%d Error: %s", resp.Payload.Error.Code, resp.Payload.Error.Message)
	}

	return nil
}

func (t Ticks) Listen(ctx context.Context, payload client.TicksFilterPayload) (<-chan tick.Tick, error) {
	ch := make(chan tick.Tick, 256)

	// Create params for channel path
	params := asyncapi.CryptellationTicksListenExchangePairParameters{
		Exchange: asyncapi.ExchangeNameSchema(payload.ExchangeName),
		Pair:     asyncapi.PairSymbolSchema(payload.PairSymbol),
	}

	// Create callback when a tick appears
	callback := func(ctx context.Context, msg asyncapi.TickMessage) {
		// Try to send tick or drop it
		select {
		case ch <- msg.ToModel():
		default:
			// Drop if it's full or closed
		}
	}

	// Listen to channel
	return ch, t.ctrl.SubscribeCryptellationTicksListenExchangePair(ctx, params, callback)
}

func (t Ticks) Unregister(ctx context.Context, payload client.TicksFilterPayload) error {
	// Generate message
	msg := asyncapi.NewRegisteringRequestMessage()
	msg.Set(payload)

	// Send message
	resp, err := t.ctrl.WaitForCryptellationTicksUnregisterResponse(ctx, &msg, func(ctx context.Context) error {
		return t.ctrl.PublishCryptellationTicksUnregisterRequest(ctx, msg)
	})
	if err != nil {
		return err
	}

	// Check error from server
	if resp.Payload.Error != nil {
		return fmt.Errorf("%d Error: %s", resp.Payload.Error.Code, resp.Payload.Error.Message)
	}

	return nil
}

func (t Ticks) Close(ctx context.Context) {
	t.ctrl.Close(ctx)
	t.broker.Close()
}
