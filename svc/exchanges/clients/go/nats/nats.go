package nats

import (
	"context"

	helpers "cryptellation/internal/asyncapi"

	asyncapi "cryptellation/svc/exchanges/api/asyncapi"
	client "cryptellation/svc/exchanges/clients/go"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	natsextension "github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"

	"cryptellation/internal/config"
	common "cryptellation/pkg/client"

	"cryptellation/svc/exchanges/pkg/exchange"
)

type nats struct {
	broker *natsextension.Controller
	ctrl   *asyncapi.UserController
	logger extensions.Logger
	name   string
}

func New(c config.NATS, options ...option) (client.Client, error) {
	var e nats

	// Execute options
	for _, option := range options {
		option(&e)
	}

	// Create a NATS Controller
	var err error
	e.broker, err = natsextension.NewController(c.URL())
	if err != nil {
		return nats{}, err
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
		return nats{}, err
	}
	e.ctrl = ctrl

	return e, nil
}

func (c nats) Read(ctx context.Context, names ...string) ([]exchange.Exchange, error) {
	// Set message
	reqMsg := asyncapi.NewListRequestMessage()
	reqMsg.Headers.ReplyTo = helpers.AddReplyToSuffix(asyncapi.ListRequestChannelPath, c.name)
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

func (c nats) ServiceInfo(ctx context.Context) (common.ServiceInfo, error) {
	// Set message
	reqMsg := asyncapi.NewServiceInfoRequestMessage()
	reqMsg.Headers.ReplyTo = helpers.AddReplyToSuffix(asyncapi.ServiceInfoRequestChannelPath, c.name)

	// Send request
	respMsg, err := c.ctrl.RequestToServiceInfoOperation(ctx, reqMsg)
	if err != nil {
		return common.ServiceInfo{}, err
	}

	return respMsg.ToModel(), nil
}

func (c nats) Close(ctx context.Context) {
	c.ctrl.Close(ctx)
	c.broker.Close()
}
