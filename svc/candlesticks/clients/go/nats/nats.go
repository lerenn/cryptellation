package nats

import (
	"context"
	"time"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
	helpers "github.com/lerenn/cryptellation/pkg/asyncapi"
	common "github.com/lerenn/cryptellation/pkg/client"
	"github.com/lerenn/cryptellation/pkg/config"
	asyncapi "github.com/lerenn/cryptellation/svc/candlesticks/api/asyncapi"
	client "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
)

type Client struct {
	broker *nats.Controller
	ctrl   *asyncapi.UserController
	logger extensions.Logger
}

type ClientOption func(c *Client)

func NewClient(c config.NATS, options ...ClientOption) (Client, error) {
	var cds Client
	var err error

	// Execute options
	for _, option := range options {
		option(&cds)
	}

	// Create a NATS Controller
	cds.broker, err = nats.NewController(c.URL())
	if err != nil {
		return Client{}, err
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
		return Client{}, err
	}
	cds.ctrl = ctrl

	return cds, nil
}

func WithLogger(logger extensions.Logger) ClientOption {
	return func(c *Client) {
		c.logger = logger
	}
}

func (c Client) Read(ctx context.Context, payload client.ReadCandlesticksPayload) (*candlestick.List, error) {
	telemetry.L(ctx).Debugf("Reading candlesticks with %+v parameters", payload)

	// Set message
	reqMsg := asyncapi.NewListRequestMessage()
	reqMsg.Headers.ReplyTo = helpers.AddReplyToSuffix(asyncapi.ListRequestChannelPath)
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

	_ = m.Loop(func(t time.Time, cs candlestick.Candlestick) (bool, error) {
		telemetry.L(ctx).Debugf("Got candlestick %s: %+v", t.Format(time.RFC3339), cs)
		return false, nil
	})

	return m, nil
}

func (c Client) ServiceInfo(ctx context.Context) (common.ServiceInfo, error) {
	// Set message
	reqMsg := asyncapi.NewServiceInfoRequestMessage()
	reqMsg.Headers.ReplyTo = helpers.AddReplyToSuffix(asyncapi.ServiceInfoRequestChannelPath)

	// Send request
	respMsg, err := c.ctrl.RequestToServiceInfoOperation(ctx, reqMsg)
	if err != nil {
		return common.ServiceInfo{}, err
	}

	return respMsg.ToModel(), nil
}

func (c Client) Close(ctx context.Context) {
	c.ctrl.Close(ctx)
	c.broker.Close()
}
