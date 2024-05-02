package nats

import (
	"context"
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
	helpers "github.com/lerenn/cryptellation/pkg/asyncapi"
	clientPkg "github.com/lerenn/cryptellation/pkg/client"
	"github.com/lerenn/cryptellation/pkg/config"
	asyncapi "github.com/lerenn/cryptellation/svc/ticks/api/asyncapi"
	client "github.com/lerenn/cryptellation/svc/ticks/clients/go"
	"github.com/lerenn/cryptellation/svc/ticks/pkg/tick"
)

type Client struct {
	broker *nats.Controller
	ctrl   *asyncapi.UserController
	logger extensions.Logger
}

type ClientOption func(t *Client)

func NewClient(c config.NATS, options ...ClientOption) (Client, error) {
	var t Client
	var err error

	// Execute options
	for _, option := range options {
		option(&t)
	}

	// Create a NATS Controller
	t.broker, err = nats.NewController(c.URL())
	if err != nil {
		return Client{}, err
	}

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
		return Client{}, err
	}
	t.ctrl = ctrl

	return t, nil
}

func WithLogger(logger extensions.Logger) ClientOption {
	return func(c *Client) {
		c.logger = logger
	}
}

func (t Client) Register(ctx context.Context, payload client.TicksFilterPayload) error {
	// Generate message
	msg := asyncapi.NewRegisteringRequestMessage()
	msg.Headers.ReplyTo = helpers.AddReplyToSuffix(asyncapi.RegisterRequestChannelPath)
	msg.Set(payload)

	// Send message
	resp, err := t.ctrl.RequestToRegisterOperation(ctx, msg)
	if err != nil {
		return err
	}

	// Check error from server
	if resp.Payload.Error != nil {
		return fmt.Errorf("%d Error: %s", resp.Payload.Error.Code, resp.Payload.Error.Message)
	}

	return nil
}

func (t Client) Listen(ctx context.Context, payload client.TicksFilterPayload) (<-chan tick.Tick, error) {
	ch := make(chan tick.Tick, 256)

	// Create params for channel path
	params := asyncapi.LiveChannelParameters{
		Exchange: payload.Exchange,
		Pair:     payload.Pair,
	}

	// Create callback when a tick appears
	callback := func(ctx context.Context, msg asyncapi.TickMessage) error {
		// Try to send tick or drop it
		select {
		case ch <- msg.ToModel():
		default:
			// Drop if it's full or closed
		}
		return nil
	}

	// Listen to channel
	return ch, t.ctrl.SubscribeToLiveOperation(ctx, params, callback)
}

func (t Client) Unregister(ctx context.Context, payload client.TicksFilterPayload) error {
	// Generate message
	msg := asyncapi.NewRegisteringRequestMessage()
	msg.Headers.ReplyTo = helpers.AddReplyToSuffix(asyncapi.UnregisterRequestChannelPath)
	msg.Set(payload)

	// Send message
	resp, err := t.ctrl.RequestToUnregisterOperation(ctx, msg)
	if err != nil {
		return err
	}

	// Check error from server
	if resp.Payload.Error != nil {
		return fmt.Errorf("%d Error: %s", resp.Payload.Error.Code, resp.Payload.Error.Message)
	}

	return nil
}

func (t Client) ServiceInfo(ctx context.Context) (clientPkg.ServiceInfo, error) {
	// Set message
	reqMsg := asyncapi.NewServiceInfoRequestMessage()
	reqMsg.Headers.ReplyTo = helpers.AddReplyToSuffix(asyncapi.ServiceInfoRequestChannelPath)

	// Send request
	respMsg, err := t.ctrl.RequestToServiceInfoOperation(ctx, reqMsg)
	if err != nil {
		return clientPkg.ServiceInfo{}, err
	}

	return respMsg.ToModel(), nil
}

func (t Client) Close(ctx context.Context) {
	t.ctrl.Close(ctx)
	t.broker.Close()
}
