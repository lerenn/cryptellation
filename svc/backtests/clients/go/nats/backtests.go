package nats

import (
	"context"
	"fmt"
	"time"

	helpers "github.com/lerenn/cryptellation/internal/asyncapi"

	common "github.com/lerenn/cryptellation/pkg/client"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/models/account"
	"github.com/lerenn/cryptellation/pkg/models/event"

	asyncapi "github.com/lerenn/cryptellation/client/api/asyncapi"
	client "github.com/lerenn/cryptellation/client/clients/go"

	"github.com/lerenn/cryptellation/ticks/pkg/tick"

	"github.com/google/uuid"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	natsextension "github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
)

type nats struct {
	broker *natsextension.Controller
	ctrl   *asyncapi.UserController

	name   string
	logger extensions.Logger
}

func New(c config.NATS, options ...option) (client.Client, error) {
	var b nats
	var err error

	// Execute options
	for _, option := range options {
		option(&b)
	}

	// Create a NATS Controller
	b.broker, err = natsextension.NewController(c.URL())
	if err != nil {
		return nil, err
	}

	// Create a logger if asked
	ctrlOpts := make([]asyncapi.ControllerOption, 0)
	if b.logger != nil {
		ctrlOpts = append(ctrlOpts, asyncapi.WithLogger(b.logger))
	} else {
		b.logger = extensions.DummyLogger{}
	}

	// Create a new user controller
	ctrl, err := asyncapi.NewUserController(b.broker, ctrlOpts...)
	if err != nil {
		return nil, err
	}
	b.ctrl = ctrl

	return &b, nil
}

func (b nats) ListenEvents(ctx context.Context, backtestID uuid.UUID) (<-chan event.Event, error) {
	ch := make(chan event.Event, 256)

	// Create callback when a tick appears
	callback := func(ctx context.Context, msg asyncapi.EventMessage) error {
		// Generate event
		e := event.Event{
			Time: time.Time(msg.Payload.Time),
			Type: event.Type(msg.Payload.Type),
		}

		// Transform message content
		switch e.Type {
		case event.TypeIsStatus:
			e.Content = event.Status{
				Finished: msg.Payload.Content.Finished,
			}
		case event.TypeIsTick:
			e.Content = tick.Tick{
				Time:     time.Time(msg.Payload.Content.Time),
				Exchange: string(msg.Payload.Content.Exchange),
				Pair:     string(msg.Payload.Content.Pair),
				Price:    msg.Payload.Content.Price,
			}
		default:
			err := fmt.Errorf("received unknown event type: %s", msg.Payload.Type)
			b.logger.Error(ctx, err.Error())
			return err
		}

		// Try to send tick or drop it
		select {
		case ch <- e:
		default:
			// Drop if it's full or closed
		}

		return nil
	}

	// Listen to channel
	return ch, b.ctrl.SubscribeToEventOperation(ctx, asyncapi.EventsChannelParameters{
		Id: backtestID.String(),
	}, callback)
}

func (b nats) Create(ctx context.Context, payload client.BacktestCreationPayload) (uuid.UUID, error) {
	// Set message
	reqMsg := asyncapi.NewCreateRequestMessage()
	reqMsg.Headers.ReplyTo = helpers.AddReplyToSuffix(asyncapi.CreateRequestChannelPath, b.name)
	reqMsg.Set(payload)

	// Send request
	respMsg, err := b.ctrl.RequestToCreateOperation(ctx, reqMsg)
	if err != nil {
		return uuid.Nil, err
	}

	// Unwrap error from message
	if err := helpers.UnwrapError(respMsg.Payload.Error); err != nil {
		return uuid.Nil, err
	}

	return uuid.Parse(respMsg.Payload.Id)
}

func (b nats) Subscribe(ctx context.Context, backtestID uuid.UUID, exchange, pair string) error {
	// Set message
	reqMsg := asyncapi.NewSubscribeRequestMessage()
	reqMsg.Headers.ReplyTo = helpers.AddReplyToSuffix(asyncapi.SubscribeRequestChannelPath, b.name)
	reqMsg.Set(backtestID, exchange, pair)

	// Send request
	respMsg, err := b.ctrl.RequestToSubscribeOperation(ctx, reqMsg)
	if err != nil {
		return err
	}

	// Unwrap error from message
	return helpers.UnwrapError(respMsg.Payload.Error)
}

func (b nats) Advance(ctx context.Context, backtestID uuid.UUID) error {
	// Set message
	reqMsg := asyncapi.NewAdvanceRequestMessage()
	reqMsg.Headers.ReplyTo = helpers.AddReplyToSuffix(asyncapi.AdvanceRequestChannelPath, b.name)
	reqMsg.Set(backtestID)

	// Send request
	respMsg, err := b.ctrl.RequestToAdvanceOperation(ctx, reqMsg)
	if err != nil {
		return err
	}

	// Unwrap error from message
	return helpers.UnwrapError(respMsg.Payload.Error)
}

func (b nats) CreateOrder(ctx context.Context, payload common.OrderCreationPayload) error {
	// Set message
	reqMsg := asyncapi.NewOrdersCreateRequestMessage()
	reqMsg.Headers.ReplyTo = helpers.AddReplyToSuffix(asyncapi.OrdersCreateRequestChannelPath, b.name)
	reqMsg.Set(payload)

	// Send request
	respMsg, err := b.ctrl.RequestToOrdersCreateOperation(ctx, reqMsg)
	if err != nil {
		return err
	}

	// Unwrap error from message
	return helpers.UnwrapError(respMsg.Payload.Error)
}

func (b nats) GetAccounts(ctx context.Context, backtestID uuid.UUID) (map[string]account.Account, error) {
	// Set message
	reqMsg := asyncapi.NewAccountsListRequestMessage()
	reqMsg.Headers.ReplyTo = helpers.AddReplyToSuffix(asyncapi.AccountsListRequestChannelPath, b.name)
	reqMsg.Set(backtestID)

	// Send request
	respMsg, err := b.ctrl.RequestToAccountsListOperation(ctx, reqMsg)
	if err != nil {
		return nil, err
	}

	// Unwrap error from message
	if err := helpers.UnwrapError(respMsg.Payload.Error); err != nil {
		return nil, err
	}

	// Convert response to model
	return respMsg.ToModel(), nil
}

func (b nats) ServiceInfo(ctx context.Context) (common.ServiceInfo, error) {
	// Set message
	reqMsg := asyncapi.NewServiceInfoRequestMessage()
	reqMsg.Headers.ReplyTo = helpers.AddReplyToSuffix(asyncapi.ServiceInfoRequestChannelPath, b.name)

	// Send request
	respMsg, err := b.ctrl.RequestToServiceInfoOperation(ctx, reqMsg)
	if err != nil {
		return common.ServiceInfo{}, err
	}

	return respMsg.ToModel(), nil
}

func (b nats) Close(ctx context.Context) {
	b.ctrl.Close(ctx)
	b.broker.Close()
}
