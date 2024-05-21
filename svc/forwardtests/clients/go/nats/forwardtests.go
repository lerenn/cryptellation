package nats

import (
	"context"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
	helpers "github.com/lerenn/cryptellation/pkg/asyncapi"
	clientPkg "github.com/lerenn/cryptellation/pkg/client"
	"github.com/lerenn/cryptellation/pkg/config"
	asyncapi "github.com/lerenn/cryptellation/svc/forwardtests/api/asyncapi"
)

type Client struct {
	broker *nats.Controller
	ctrl   *asyncapi.UserController
	logger extensions.Logger
}

type ForwardtestsOption func(b *Client)

func NewClient(c config.NATS, options ...ForwardtestsOption) (Client, error) {
	var b Client
	var err error

	// Execute options
	for _, option := range options {
		option(&b)
	}

	// Create a NATS Controller
	b.broker, err = nats.NewController(c.URL())
	if err != nil {
		return Client{}, err
	}

	// Create a logger if asked
	ctrlOpts := make([]asyncapi.ControllerOption, 0)
	if b.logger != nil {
		ctrlOpts = append(ctrlOpts, asyncapi.WithLogger(b.logger))
	} else {
		b.logger = extensions.DummyLogger{}
	}

	// Create a new user controller
	ctrl, err := asyncapi.NewUserController(b.broker, ctrlOpts...)
	if err != nil {
		return Client{}, err
	}
	b.ctrl = ctrl

	return b, nil
}

func WithLogger(logger extensions.Logger) ForwardtestsOption {
	return func(b *Client) {
		b.logger = logger
	}
}

func (b Client) ServiceInfo(ctx context.Context) (clientPkg.ServiceInfo, error) {
	// Set message
	reqMsg := asyncapi.NewServiceInfoRequestMessage()
	reqMsg.Headers.ReplyTo = helpers.AddReplyToSuffix(asyncapi.ServiceInfoRequestChannelPath)

	// Send request
	respMsg, err := b.ctrl.RequestToServiceInfoOperation(ctx, reqMsg)
	if err != nil {
		return clientPkg.ServiceInfo{}, err
	}

	return respMsg.ToModel(), nil
}

func (b Client) Close(ctx context.Context) {
	b.ctrl.Close(ctx)
	b.broker.Close()
}
