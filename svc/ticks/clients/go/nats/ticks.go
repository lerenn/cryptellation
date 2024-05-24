package nats

import (
	"context"
	"time"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
	helpers "github.com/lerenn/cryptellation/pkg/asyncapi"
	common "github.com/lerenn/cryptellation/pkg/client"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/models/event"
	asyncapi "github.com/lerenn/cryptellation/svc/ticks/api/asyncapi"
	"github.com/lerenn/cryptellation/svc/ticks/pkg/tick"
)

type Client struct {
	broker *nats.Controller
	ctrl   *asyncapi.UserController
	logger extensions.Logger
}

type ClientOption func(t *Client)

func NewClient(c config.NATS, options ...ClientOption) (Client, error) {
	var t Client
	var err error

	// Execute options
	for _, option := range options {
		option(&t)
	}

	// Create a NATS Controller
	t.broker, err = nats.NewController(c.URL())
	if err != nil {
		return Client{}, err
	}

	// Create a logger if asked
	ctrlOpts := make([]asyncapi.ControllerOption, 0)
	if t.logger != nil {
		ctrlOpts = append(ctrlOpts, asyncapi.WithLogger(t.logger))
	} else {
		t.logger = extensions.DummyLogger{}
	}

	// Create a new user controller
	ctrl, err := asyncapi.NewUserController(t.broker, ctrlOpts...)
	if err != nil {
		return Client{}, err
	}
	t.ctrl = ctrl

	return t, nil
}

func WithLogger(logger extensions.Logger) ClientOption {
	return func(c *Client) {
		c.logger = logger
	}
}

func (t Client) sendSubscriptionRequest(ctx context.Context, sub event.TickSubscription) error {
	// Create message
	msg := asyncapi.NewListeningNotificationMessage()
	msg.FromModel(sub)

	// Send message
	return t.ctrl.SendToListeningOperation(ctx, msg)
}

func (t Client) SubscribeToTicks(ctx context.Context, sub event.TickSubscription) (<-chan tick.Tick, error) {
	// Send subscription request periodically
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(1 * time.Second):
				if err := t.sendSubscriptionRequest(ctx, sub); err != nil {
					// Log error
					continue
				}
			}
		}
	}()

	// Subscribe to new ticks
	ch := make(chan tick.Tick, 16)
	if err := t.ctrl.SubscribeToSendNewTicksOperation(ctx, asyncapi.LiveChannelParameters{
		Exchange: sub.Exchange,
		Pair:     sub.Pair,
	}, func(ctx context.Context, msg asyncapi.TickMessage) error {
		ch <- msg.ToModel()
		return nil
	}); err != nil {
		return nil, err
	}

	return ch, nil
}

func (t Client) ServiceInfo(ctx context.Context) (common.ServiceInfo, error) {
	// Set message
	reqMsg := asyncapi.NewServiceInfoRequestMessage()
	reqMsg.Headers.ReplyTo = helpers.AddReplyToSuffix(asyncapi.ServiceInfoRequestChannelPath)

	// Send request
	respMsg, err := t.ctrl.RequestToServiceInfoOperation(ctx, reqMsg)
	if err != nil {
		return common.ServiceInfo{}, err
	}

	return respMsg.ToModel(), nil
}

func (t Client) Close(ctx context.Context) {
	t.ctrl.Close(ctx)
	t.broker.Close()
}
