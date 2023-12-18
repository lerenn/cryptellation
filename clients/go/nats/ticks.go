package nats

import (
	"context"
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
	asyncapi "github.com/lerenn/cryptellation/api/asyncapi/ticks"
	client "github.com/lerenn/cryptellation/clients/go"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/models/tick"
)

type Ticks struct {
	broker *nats.Controller
	ctrl   *asyncapi.UserController
	logger extensions.Logger
}

type TicksOption func(t *Ticks)

func NewTicks(c config.NATS, options ...TicksOption) (client.Ticks, error) {
	var t Ticks

	// Execute options
	for _, option := range options {
		option(&t)
	}

	// Create a NATS Controller
	t.broker = nats.NewController(c.URL())

	// Create a logger if asked
	ctrlOpts := make([]asyncapi.ControllerOption, 0)
	if t.logger != nil {
		ctrlOpts = append(ctrlOpts, asyncapi.WithLogger(t.logger))
	} else {
		t.logger = extensions.DummyLogger{}
	}

	// Create a new user controller
	ctrl, err := asyncapi.NewUserController(t.broker, ctrlOpts...)
	if err != nil {
		return nil, err
	}
	t.ctrl = ctrl

	return t, nil
}

func WithTicksLogger(logger extensions.Logger) TicksOption {
	return func(c *Ticks) {
		c.logger = logger
	}
}

func (t Ticks) Register(ctx context.Context, payload client.TicksFilterPayload) error {
	// Generate message
	msg := asyncapi.NewRegisteringRequestMessage()
	msg.Set(payload)

	// Send message
	resp, err := t.ctrl.WaitForRegisterToTicksResponse(ctx, &msg, func(ctx context.Context) error {
		return t.ctrl.PublishRegisterToTicksRequest(ctx, msg)
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
	params := asyncapi.CryptellationTicksLiveParameters{
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
	return ch, t.ctrl.SubscribeWatchTicks(ctx, params, callback)
}

func (t Ticks) Unregister(ctx context.Context, payload client.TicksFilterPayload) error {
	// Generate message
	msg := asyncapi.NewRegisteringRequestMessage()
	msg.Set(payload)

	// Send message
	resp, err := t.ctrl.WaitForUnregisterToTicksResponse(ctx, &msg, func(ctx context.Context) error {
		return t.ctrl.PublishUnregisterToTicksRequest(ctx, msg)
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

func (t Ticks) ServiceInfo(ctx context.Context) (client.ServiceInfo, error) {
	// Set message
	reqMsg := asyncapi.NewServiceInfoRequestMessage()

	// Send request
	respMsg, err := t.ctrl.WaitForServiceInfoResponse(ctx, &reqMsg, func(ctx context.Context) error {
		return t.ctrl.PublishServiceInfoRequest(ctx, reqMsg)
	})
	if err != nil {
		return client.ServiceInfo{}, err
	}

	return respMsg.ToModel(), nil
}

func (t Ticks) Close(ctx context.Context) {
	t.ctrl.Close(ctx)
	t.broker.Close()
}
