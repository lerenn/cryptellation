package nats

import (
	"context"
	"time"

	helpers "cryptellation/internal/asyncapi"
	"cryptellation/pkg/adapters/telemetry"
	common "cryptellation/pkg/client"
	"cryptellation/pkg/config"

	asyncapi "cryptellation/svc/candlesticks/api/asyncapi"
	client "cryptellation/svc/candlesticks/clients/go"
	"cryptellation/svc/candlesticks/pkg/candlestick"

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
	var cds nats
	var err error

	// Execute options
	for _, option := range options {
		option(&cds)
	}

	// Create a NATS Controller
	cds.broker, err = natsextension.NewController(c.URL())
	if err != nil {
		return nats{}, err
	}

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
		return nats{}, err
	}
	cds.ctrl = ctrl

	return cds, nil
}

func (c nats) Read(ctx context.Context, payload client.ReadCandlesticksPayload) (*candlestick.List, error) {
	telemetry.L(ctx).Debugf("Reading candlesticks with %+v parameters", payload)

	// Set message
	reqMsg := asyncapi.NewListRequestMessage()
	reqMsg.Headers.ReplyTo = helpers.AddReplyToSuffix(asyncapi.ListRequestChannelPath, c.name)
	reqMsg.Set(payload)

	// Send request
	respMsg, err := c.ctrl.RequestToListOperation(ctx, reqMsg)
	if err != nil {
		return nil, err
	}

	// Unwrap error from message
	if err := helpers.UnwrapError(respMsg.Payload.Error); err != nil {
		return nil, err
	}

	// To candlestick list
	m, err := respMsg.ToModel(payload.Exchange, payload.Pair, payload.Period)
	if err != nil {
		return nil, err
	}

	// Debug content
	_ = m.Loop(func(t time.Time, cs candlestick.Candlestick) (bool, error) {
		telemetry.L(ctx).Debugf("Got candlestick %s: %+v", t.Format(time.RFC3339), cs)
		return false, nil
	})

	return m, nil
}

func (c nats) ServiceInfo(ctx context.Context) (common.ServiceInfo, error) {
	// Set message
	reqMsg := asyncapi.NewServiceInfoRequestMessage()
	reqMsg.Headers.ReplyTo = helpers.AddReplyToSuffix(asyncapi.ServiceInfoRequestChannelPath, c.name)

	// Send request
	respMsg, err := c.ctrl.RequestToServiceInfoOperation(ctx, reqMsg)
	if err != nil {
		return common.ServiceInfo{}, err
	}

	return respMsg.ToModel(), nil
}

func (c nats) Close(ctx context.Context) {
	c.ctrl.Close(ctx)
	c.broker.Close()
}
