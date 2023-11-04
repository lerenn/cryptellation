package nats

import (
	"context"
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/loggers"
	client "github.com/lerenn/cryptellation/clients/go"
	asyncapi "github.com/lerenn/cryptellation/pkg/asyncapi/candlesticks"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/models/candlestick"
)

type Candlesticks struct {
	broker *nats.Controller
	ctrl   *asyncapi.UserController
	logger extensions.Logger
}

func NewCandlesticks(c config.NATS) (client.Candlesticks, error) {
	// Create a NATS Controller
	broker := nats.NewController(c.URL())

	// Create a logger
	logger := loggers.NewECS()

	// Create a new user controller
	ctrl, err := asyncapi.NewUserController(broker, asyncapi.WithLogger(logger))
	if err != nil {
		return nil, err
	}

	return Candlesticks{
		broker: broker,
		ctrl:   ctrl,
		logger: logger,
	}, nil
}

func (c Candlesticks) Read(ctx context.Context, payload client.ReadCandlesticksPayload) (*candlestick.List, error) {
	// Set message
	reqMsg := asyncapi.NewCandlesticksListRequestMessage()
	reqMsg.Set(payload)

	// Send request
	respMsg, err := c.ctrl.WaitForCryptellationCandlesticksListResponse(ctx, &reqMsg, func(ctx context.Context) error {
		return c.ctrl.PublishCryptellationCandlesticksListRequest(ctx, reqMsg)
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

func (c Candlesticks) Close(ctx context.Context) {
	c.ctrl.Close(ctx)
	c.broker.Close()
}
