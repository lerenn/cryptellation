package nats

import (
	"context"
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/loggers"
	client "github.com/lerenn/cryptellation/clients/go"
	asyncapi "github.com/lerenn/cryptellation/pkg/asyncapi/indicators"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/models/timeserie"
)

type Indicators struct {
	broker *nats.Controller
	ctrl   *asyncapi.UserController
	logger extensions.Logger
}

func NewIndicators(c config.NATS) (client.Indicators, error) {
	// Create a NATS Controller
	broker := nats.NewController(c.URL())

	// Create a logger
	logger := loggers.NewECS()

	// Create a new user controller
	ctrl, err := asyncapi.NewUserController(broker, asyncapi.WithLogger(logger))
	if err != nil {
		return nil, err
	}

	return Indicators{
		broker: broker,
		ctrl:   ctrl,
		logger: logger,
	}, nil
}

func (ex Indicators) SMA(ctx context.Context, payload client.SMAPayload) (*timeserie.TimeSerie[float64], error) {
	// Set message
	reqMsg := asyncapi.NewSmaRequestMessage()
	reqMsg.Set(payload)

	// Send request
	respMsg, err := ex.ctrl.WaitForCryptellationIndicatorsSmaResponse(ctx, &reqMsg, func(ctx context.Context) error {
		return ex.ctrl.PublishCryptellationIndicatorsSmaRequest(ctx, reqMsg)
	})
	if err != nil {
		return nil, err
	}

	// Check error
	if respMsg.Payload.Error != nil {
		return nil, fmt.Errorf("%d Error: %s", respMsg.Payload.Error.Code, respMsg.Payload.Error.Message)
	}

	// To indicator list
	return respMsg.ToModel(), nil
}

func (ex Indicators) Close(ctx context.Context) {
	ex.ctrl.Close(ctx)
	ex.broker.Close()
}
