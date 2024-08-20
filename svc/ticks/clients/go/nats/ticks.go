package nats

import (
	"context"
	"time"

	helpers "cryptellation/internal/asyncapi"
	common "cryptellation/pkg/client"
	"cryptellation/pkg/config"
	"cryptellation/pkg/models/event"

	asyncapi "cryptellation/svc/ticks/api/asyncapi"
	client "cryptellation/svc/ticks/clients/go"
	"cryptellation/svc/ticks/pkg/tick"

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
	var t nats
	var err error

	// Execute options
	for _, option := range options {
		option(&t)
	}

	// Create a NATS Controller
	t.broker, err = natsextension.NewController(c.URL())
	if err != nil {
		return nats{}, err
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
		return nats{}, err
	}
	t.ctrl = ctrl

	return t, nil
}

func (t nats) sendSubscriptionRequest(ctx context.Context, sub event.TickSubscription) error {
	// Create message
	msg := asyncapi.NewListeningNotificationMessage()
	msg.FromModel(sub)

	// Send message
	return t.ctrl.SendToListeningOperation(ctx, msg)
}

func (t nats) SubscribeToTicks(ctx context.Context, sub event.TickSubscription) (<-chan tick.Tick, error) {
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

func (t nats) ServiceInfo(ctx context.Context) (common.ServiceInfo, error) {
	// Set message
	reqMsg := asyncapi.NewServiceInfoRequestMessage()
	reqMsg.Headers.ReplyTo = helpers.AddReplyToSuffix(asyncapi.ServiceInfoRequestChannelPath, t.name)

	// Send request
	respMsg, err := t.ctrl.RequestToServiceInfoOperation(ctx, reqMsg)
	if err != nil {
		return common.ServiceInfo{}, err
	}

	return respMsg.ToModel(), nil
}

func (t nats) Close(ctx context.Context) {
	t.ctrl.Close(ctx)
	t.broker.Close()
}
