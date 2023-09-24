package nats

import (
	"context"
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/loggers"
	client "github.com/lerenn/cryptellation/clients/go"
	"github.com/lerenn/cryptellation/internal/ctrl/exchanges/events"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/models/exchange"
)

type Exchanges struct {
	broker *nats.Controller
	ctrl   *events.UserController
	logger extensions.Logger
}

func NewExchanges(c config.NATS) (client.Exchanges, error) {
	// Create a NATS Controller
	broker := nats.NewController(c.URL())

	// Create a logger
	logger := loggers.NewECS()

	// Create a new user controller
	ctrl, err := events.NewUserController(broker, events.WithLogger(logger))
	if err != nil {
		return nil, err
	}

	return Exchanges{
		broker: broker,
		ctrl:   ctrl,
		logger: logger,
	}, nil
}

func (ex Exchanges) Read(ctx context.Context, names ...string) ([]exchange.Exchange, error) {
	// Set message
	reqMsg := events.NewExchangesRequestMessage()
	reqMsg.Set(names...)

	// Send request
	respMsg, err := ex.ctrl.WaitForCryptellationExchangesListResponse(ctx, &reqMsg, func(ctx context.Context) error {
		return ex.ctrl.PublishCryptellationExchangesListRequest(ctx, reqMsg)
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

func (ex Exchanges) Close(ctx context.Context) {
	ex.ctrl.Close(ctx)
	ex.broker.Close()
}
