package nats

import (
	"context"
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
	clientPkg "github.com/lerenn/cryptellation/pkg/client"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/timeserie"
	asyncapi "github.com/lerenn/cryptellation/svc/indicators/api/asyncapi"
	client "github.com/lerenn/cryptellation/svc/indicators/clients/go"
)

type Client struct {
	broker *nats.Controller
	ctrl   *asyncapi.UserController
	logger extensions.Logger
}

type IndicatorsOption func(i *Client)

func NewClient(c config.NATS, options ...IndicatorsOption) (Client, error) {
	var i Client
	var err error

	// Execute options
	for _, option := range options {
		option(&i)
	}

	// Create a NATS Controller
	i.broker, err = nats.NewController(c.URL())
	if err != nil {
		return Client{}, err
	}

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
		return Client{}, err
	}
	i.ctrl = ctrl

	return i, nil
}

func WithIndicatorsLogger(logger extensions.Logger) IndicatorsOption {
	return func(c *Client) {
		c.logger = logger
	}
}

func (ids Client) SMA(ctx context.Context, payload client.SMAPayload) (*timeserie.TimeSerie[float64], error) {
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

func (ids Client) ServiceInfo(ctx context.Context) (clientPkg.ServiceInfo, error) {
	// Set message
	reqMsg := asyncapi.NewServiceInfoRequestMessage()

	// Send request
	respMsg, err := ids.ctrl.WaitForServiceInfoResponse(ctx, &reqMsg, func(ctx context.Context) error {
		return ids.ctrl.PublishServiceInfoRequest(ctx, reqMsg)
	})
	if err != nil {
		return clientPkg.ServiceInfo{}, err
	}

	return respMsg.ToModel(), nil
}

func (ids Client) Close(ctx context.Context) {
	ids.ctrl.Close(ctx)
	ids.broker.Close()
}
