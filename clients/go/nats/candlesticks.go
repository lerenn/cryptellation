package nats

import (
	"context"
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
	asyncapi "github.com/lerenn/cryptellation/api/asyncapi/candlesticks"
	client "github.com/lerenn/cryptellation/clients/go"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/models/candlestick"
)

type Candlesticks struct {
	broker *nats.Controller
	ctrl   *asyncapi.UserController
	logger extensions.Logger
}

type CandlesticksOption func(c *Candlesticks)

func NewCandlesticks(c config.NATS, options ...CandlesticksOption) (client.Candlesticks, error) {
	var cds Candlesticks

	// Execute options
	for _, option := range options {
		option(&cds)
	}

	// Create a NATS Controller
	cds.broker = nats.NewController(c.URL())

	// Create a logger if asked
	ctrlOpts := make([]asyncapi.ControllerOption, 0)
	if cds.logger != nil {
		ctrlOpts = append(ctrlOpts, asyncapi.WithLogger(cds.logger))
	} else {
		cds.logger = extensions.DummyLogger{}
	}

	// Create a new user controller
	ctrl, err := asyncapi.NewUserController(cds.broker, ctrlOpts...)
	if err != nil {
		return nil, err
	}
	cds.ctrl = ctrl

	return cds, nil
}

func WithCandlesticksLogger(logger extensions.Logger) CandlesticksOption {
	return func(c *Candlesticks) {
		c.logger = logger
	}
}

func (c Candlesticks) Read(ctx context.Context, payload client.ReadCandlesticksPayload) (*candlestick.List, error) {
	// Set message
	reqMsg := asyncapi.NewListCandlesticksRequestMessage()
	reqMsg.Set(payload)

	// Send request
	respMsg, err := c.ctrl.WaitForListCandlesticksResponse(ctx, &reqMsg, func(ctx context.Context) error {
		return c.ctrl.PublishListCandlesticksRequest(ctx, reqMsg)
	})
	if err != nil {
		return nil, err
	}

	// Check error
	if respMsg.Payload.Error != nil {
		return nil, fmt.Errorf("%d Error: %s", respMsg.Payload.Error.Code, respMsg.Payload.Error.Message)
	}

	// To candlestick list
	return respMsg.ToModel(payload.ExchangeName, payload.PairSymbol, payload.Period)
}

func (c Candlesticks) ServiceInfo(ctx context.Context) (client.ServiceInfo, error) {
	// Set message
	reqMsg := asyncapi.NewServiceInfoRequestMessage()

	// Send request
	respMsg, err := c.ctrl.WaitForServiceInfoResponse(ctx, &reqMsg, func(ctx context.Context) error {
		return c.ctrl.PublishServiceInfoRequest(ctx, reqMsg)
	})
	if err != nil {
		return client.ServiceInfo{}, err
	}

	return respMsg.ToModel(), nil
}

func (c Candlesticks) Close(ctx context.Context) {
	c.ctrl.Close(ctx)
	c.broker.Close()
}
