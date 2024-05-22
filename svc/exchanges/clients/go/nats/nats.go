package nats

import (
	"context"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
	helpers "github.com/lerenn/cryptellation/pkg/asyncapi"
	asyncapi "github.com/lerenn/cryptellation/svc/exchanges/api/asyncapi"

	"github.com/lerenn/cryptellation/pkg/client"
	"github.com/lerenn/cryptellation/pkg/config"

	"github.com/lerenn/cryptellation/svc/exchanges/pkg/exchange"
)

type Client struct {
	broker *nats.Controller
	ctrl   *asyncapi.UserController
	logger extensions.Logger
}

type ClientOption func(e *Client)

func NewClient(c config.NATS, options ...ClientOption) (Client, error) {
	var e Client

	// Execute options
	for _, option := range options {
		option(&e)
	}

	// Create a NATS Controller
	var err error
	e.broker, err = nats.NewController(c.URL())
	if err != nil {
		return Client{}, err
	}

	// Create a logger if asked
	ctrlOpts := make([]asyncapi.ControllerOption, 0)
	if e.logger != nil {
		ctrlOpts = append(ctrlOpts, asyncapi.WithLogger(e.logger))
	} else {
		e.logger = extensions.DummyLogger{}
	}

	// Create a new user controller
	ctrl, err := asyncapi.NewUserController(e.broker, ctrlOpts...)
	if err != nil {
		return Client{}, err
	}
	e.ctrl = ctrl

	return e, nil
}

func WithLogger(logger extensions.Logger) ClientOption {
	return func(c *Client) {
		c.logger = logger
	}
}

func (c Client) Read(ctx context.Context, names ...string) ([]exchange.Exchange, error) {
	// Set message
	reqMsg := asyncapi.NewListRequestMessage()
	reqMsg.Headers.ReplyTo = helpers.AddReplyToSuffix(asyncapi.ListRequestChannelPath)
	reqMsg.Set(names...)

	// Send request
	respMsg, err := c.ctrl.RequestToListOperation(ctx, reqMsg)
	if err != nil {
		return nil, err
	}

	// Unwrap error from message
	if err := helpers.UnwrapError(respMsg.Payload.Error); err != nil {
		return nil, err
	}

	// To exchange list
	return respMsg.ToModel(), nil
}

func (c Client) ServiceInfo(ctx context.Context) (client.ServiceInfo, error) {
	// Set message
	reqMsg := asyncapi.NewServiceInfoRequestMessage()
	reqMsg.Headers.ReplyTo = helpers.AddReplyToSuffix(asyncapi.ServiceInfoRequestChannelPath)

	// Send request
	respMsg, err := c.ctrl.RequestToServiceInfoOperation(ctx, reqMsg)
	if err != nil {
		return client.ServiceInfo{}, err
	}

	return respMsg.ToModel(), nil
}

func (c Client) Close(ctx context.Context) {
	c.ctrl.Close(ctx)
	c.broker.Close()
}
