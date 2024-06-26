package nats

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
	helpers "github.com/lerenn/cryptellation/pkg/asyncapi"
	common "github.com/lerenn/cryptellation/pkg/client"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/models/account"
	"github.com/lerenn/cryptellation/pkg/models/event"
	asyncapi "github.com/lerenn/cryptellation/svc/backtests/api/asyncapi"
	client "github.com/lerenn/cryptellation/svc/backtests/clients/go"
	"github.com/lerenn/cryptellation/svc/ticks/pkg/tick"
)

type Client struct {
	broker *nats.Controller
	ctrl   *asyncapi.UserController

	name   string
	logger extensions.Logger
}

type BacktestsOption func(b *Client)

func NewClient(c config.NATS, options ...BacktestsOption) (Client, error) {
	var b Client
	var err error

	// Execute options
	for _, option := range options {
		option(&b)
	}

	// Create a NATS Controller
	b.broker, err = nats.NewController(c.URL())
	if err != nil {
		return Client{}, err
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
		return Client{}, err
	}
	b.ctrl = ctrl

	return b, nil
}

func WithLogger(logger extensions.Logger) BacktestsOption {
	return func(b *Client) {
		b.logger = logger
	}
}

func WithName(name string) BacktestsOption {
	return func(b *Client) {
		b.name = name
	}
}

func (b Client) ListenEvents(ctx context.Context, backtestID uuid.UUID) (<-chan event.Event, error) {
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

func (b Client) Create(ctx context.Context, payload client.BacktestCreationPayload) (uuid.UUID, error) {
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

func (b Client) Subscribe(ctx context.Context, backtestID uuid.UUID, exchange, pair string) error {
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

func (b Client) Advance(ctx context.Context, backtestID uuid.UUID) error {
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

func (b Client) CreateOrder(ctx context.Context, payload common.OrderCreationPayload) error {
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

func (b Client) GetAccounts(ctx context.Context, backtestID uuid.UUID) (map[string]account.Account, error) {
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

func (b Client) ServiceInfo(ctx context.Context) (common.ServiceInfo, error) {
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

func (b Client) Close(ctx context.Context) {
	b.ctrl.Close(ctx)
	b.broker.Close()
}
