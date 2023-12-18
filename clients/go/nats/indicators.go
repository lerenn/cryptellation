package nats

import (
	"context"
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
	asyncapi "github.com/lerenn/cryptellation/api/asyncapi/indicators"
	client "github.com/lerenn/cryptellation/clients/go"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/models/timeserie"
)

type Indicators struct {
	broker *nats.Controller
	ctrl   *asyncapi.UserController
	logger extensions.Logger
}

type IndicatorsOption func(i *Indicators)

func NewIndicators(c config.NATS, options ...IndicatorsOption) (client.Indicators, error) {
	var i Indicators

	// Execute options
	for _, option := range options {
		option(&i)
	}

	// Create a NATS Controller
	i.broker = nats.NewController(c.URL())

	// Create a logger if asked
	ctrlOpts := make([]asyncapi.ControllerOption, 0)
	if i.logger != nil {
		ctrlOpts = append(ctrlOpts, asyncapi.WithLogger(i.logger))
	} else {
		i.logger = extensions.DummyLogger{}
	}

	// Create a new user controller
	ctrl, err := asyncapi.NewUserController(i.broker, ctrlOpts...)
	if err != nil {
		return nil, err
	}
	i.ctrl = ctrl

	return i, nil
}

func WithIndicatorsLogger(logger extensions.Logger) IndicatorsOption {
	return func(c *Indicators) {
		c.logger = logger
	}
}

func (ids Indicators) SMA(ctx context.Context, payload client.SMAPayload) (*timeserie.TimeSerie[float64], error) {
	// Set message
	reqMsg := asyncapi.NewGetSMARequestMessage()
	reqMsg.Set(payload)

	// Send request
	respMsg, err := ids.ctrl.WaitForGetSMAResponse(ctx, &reqMsg, func(ctx context.Context) error {
		return ids.ctrl.PublishGetSMARequest(ctx, reqMsg)
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

func (ids Indicators) ServiceInfo(ctx context.Context) (client.ServiceInfo, error) {
	// Set message
	reqMsg := asyncapi.NewServiceInfoRequestMessage()

	// Send request
	respMsg, err := ids.ctrl.WaitForServiceInfoResponse(ctx, &reqMsg, func(ctx context.Context) error {
		return ids.ctrl.PublishServiceInfoRequest(ctx, reqMsg)
	})
	if err != nil {
		return client.ServiceInfo{}, err
	}

	return respMsg.ToModel(), nil
}

func (ids Indicators) Close(ctx context.Context) {
	ids.ctrl.Close(ctx)
	ids.broker.Close()
}
