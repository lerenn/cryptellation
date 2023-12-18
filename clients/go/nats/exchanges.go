package nats

import (
	"context"
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
	asyncapi "github.com/lerenn/cryptellation/api/asyncapi/exchanges"
	client "github.com/lerenn/cryptellation/clients/go"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/models/exchange"
)

type Exchanges struct {
	broker *nats.Controller
	ctrl   *asyncapi.UserController
	logger extensions.Logger
}

type ExchangesOption func(e *Exchanges)

func NewExchanges(c config.NATS, options ...ExchangesOption) (client.Exchanges, error) {
	var e Exchanges

	// Execute options
	for _, option := range options {
		option(&e)
	}

	// Create a NATS Controller
	e.broker = nats.NewController(c.URL())

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
		return nil, err
	}
	e.ctrl = ctrl

	return e, nil
}

func WithExchangesLogger(logger extensions.Logger) ExchangesOption {
	return func(c *Exchanges) {
		c.logger = logger
	}
}

func (ex Exchanges) Read(ctx context.Context, names ...string) ([]exchange.Exchange, error) {
	// Set message
	reqMsg := asyncapi.NewListExchangesRequestMessage()
	reqMsg.Set(names...)

	// Send request
	respMsg, err := ex.ctrl.WaitForListExchangesResponse(ctx, &reqMsg, func(ctx context.Context) error {
		return ex.ctrl.PublishListExchangesRequest(ctx, reqMsg)
	})
	if err != nil {
		return nil, err
	}

	// Check error
	if respMsg.Payload.Error != nil {
		return nil, fmt.Errorf("%d Error: %s", respMsg.Payload.Error.Code, respMsg.Payload.Error.Message)
	}

	// To exchange list
	return respMsg.ToModel(), nil
}

func (ex Exchanges) ServiceInfo(ctx context.Context) (client.ServiceInfo, error) {
	// Set message
	reqMsg := asyncapi.NewServiceInfoRequestMessage()

	// Send request
	respMsg, err := ex.ctrl.WaitForServiceInfoResponse(ctx, &reqMsg, func(ctx context.Context) error {
		return ex.ctrl.PublishServiceInfoRequest(ctx, reqMsg)
	})
	if err != nil {
		return client.ServiceInfo{}, err
	}

	return respMsg.ToModel(), nil
}

func (ex Exchanges) Close(ctx context.Context) {
	ex.ctrl.Close(ctx)
	ex.broker.Close()
}
