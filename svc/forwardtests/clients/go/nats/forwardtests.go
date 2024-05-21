package nats

import (
	"context"

	"github.com/google/uuid"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
	helpers "github.com/lerenn/cryptellation/pkg/asyncapi"
	clientPkg "github.com/lerenn/cryptellation/pkg/client"
	"github.com/lerenn/cryptellation/pkg/config"
	asyncapi "github.com/lerenn/cryptellation/svc/forwardtests/api/asyncapi"
	"github.com/lerenn/cryptellation/svc/forwardtests/pkg/forwardtest"
)

type Client struct {
	broker *nats.Controller
	ctrl   *asyncapi.UserController
	logger extensions.Logger
}

type ForwardTestsOption func(b *Client)

func NewClient(c config.NATS, options ...ForwardTestsOption) (Client, error) {
	var cl Client
	var err error

	// Execute options
	for _, option := range options {
		option(&cl)
	}

	// Create a NATS Controller
	cl.broker, err = nats.NewController(c.URL())
	if err != nil {
		return Client{}, err
	}

	// Create a logger if asked
	ctrlOpts := make([]asyncapi.ControllerOption, 0)
	if cl.logger != nil {
		ctrlOpts = append(ctrlOpts, asyncapi.WithLogger(cl.logger))
	} else {
		cl.logger = extensions.DummyLogger{}
	}

	// Create a new user controller
	ctrl, err := asyncapi.NewUserController(cl.broker, ctrlOpts...)
	if err != nil {
		return Client{}, err
	}
	cl.ctrl = ctrl

	return cl, nil
}

func (cl Client) CreateForwardTest(ctx context.Context, payload forwardtest.NewPayload) (uuid.UUID, error) {
	// Set message
	reqMsg := asyncapi.NewCreateRequestMessage()
	reqMsg.Set(payload)

	// Send request
	respMsg, err := cl.ctrl.RequestToCreateOperation(ctx, reqMsg)
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.Parse(respMsg.Payload.Id)

}

func WithLogger(logger extensions.Logger) ForwardTestsOption {
	return func(b *Client) {
		b.logger = logger
	}
}

func (cl Client) ServiceInfo(ctx context.Context) (clientPkg.ServiceInfo, error) {
	// Set message
	reqMsg := asyncapi.NewServiceInfoRequestMessage()
	reqMsg.Headers.ReplyTo = helpers.AddReplyToSuffix(asyncapi.ServiceInfoRequestChannelPath)

	// Send request
	respMsg, err := cl.ctrl.RequestToServiceInfoOperation(ctx, reqMsg)
	if err != nil {
		return clientPkg.ServiceInfo{}, err
	}

	return respMsg.ToModel(), nil
}

func (cl Client) Close(ctx context.Context) {
	cl.ctrl.Close(ctx)
	cl.broker.Close()
}
