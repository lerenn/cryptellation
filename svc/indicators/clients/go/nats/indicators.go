package nats

import (
	"context"

	helpers "cryptellation/internal/asyncapi"
	"cryptellation/internal/config"
	common "cryptellation/pkg/client"
	"cryptellation/pkg/models/timeserie"

	asyncapi "cryptellation/svc/indicators/api/asyncapi"
	client "cryptellation/svc/indicators/clients/go"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	natsextension "github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
)

type nats struct {
	broker *natsextension.Controller
	ctrl   *asyncapi.UserController
	logger extensions.Logger
	name   string
}

func New(c config.NATS, options ...option) (client.Client, error) {
	var i nats
	var err error

	// Execute options
	for _, option := range options {
		option(&i)
	}

	// Create a NATS Controller
	i.broker, err = natsextension.NewController(c.URL())
	if err != nil {
		return nats{}, err
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
		return nats{}, err
	}
	i.ctrl = ctrl

	return i, nil
}

func (ids nats) SMA(ctx context.Context, payload client.SMAPayload) (*timeserie.TimeSerie[float64], error) {
	// Set message
	reqMsg := asyncapi.NewSMARequestMessage()
	reqMsg.Headers.ReplyTo = helpers.AddReplyToSuffix(asyncapi.SMARequestChannelPath, ids.name)
	reqMsg.Set(payload)

	// Send request
	respMsg, err := ids.ctrl.RequestToSMAOperation(ctx, reqMsg)
	if err != nil {
		return nil, err
	}

	// Unwrap error from message
	if err := helpers.UnwrapError(respMsg.Payload.Error); err != nil {
		return nil, err
	}

	// To indicator list
	return respMsg.ToModel(), nil
}

func (ids nats) ServiceInfo(ctx context.Context) (common.ServiceInfo, error) {
	// Set message
	reqMsg := asyncapi.NewServiceInfoRequestMessage()
	reqMsg.Headers.ReplyTo = helpers.AddReplyToSuffix(asyncapi.ServiceInfoRequestChannelPath, ids.name)

	// Send request
	respMsg, err := ids.ctrl.RequestToServiceInfoOperation(ctx, reqMsg)
	if err != nil {
		return common.ServiceInfo{}, err
	}

	return respMsg.ToModel(), nil
}

func (ids nats) Close(ctx context.Context) {
	ids.ctrl.Close(ctx)
	ids.broker.Close()
}
