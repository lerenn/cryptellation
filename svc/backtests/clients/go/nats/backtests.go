package nats

import (
	"context"
	"fmt"
	"time"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
	helpers "github.com/lerenn/cryptellation/pkg/asyncapi"
	clientPkg "github.com/lerenn/cryptellation/pkg/client"
	"github.com/lerenn/cryptellation/pkg/config"
	asyncapi "github.com/lerenn/cryptellation/svc/backtests/api/asyncapi"
	client "github.com/lerenn/cryptellation/svc/backtests/clients/go"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/account"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/event"
	"github.com/lerenn/cryptellation/svc/ticks/pkg/tick"
)

type Client struct {
	broker *nats.Controller
	ctrl   *asyncapi.UserController

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

func (b Client) ListenEvents(ctx context.Context, backtestID uint) (<-chan event.Event, error) {
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
		Id: fmt.Sprintf("%d", backtestID),
	}, callback)
}

func (b Client) Create(ctx context.Context, payload client.BacktestCreationPayload) (uint, error) {
	// Set message
	reqMsg := asyncapi.NewCreateRequestMessage()
	reqMsg.Headers.ReplyTo = helpers.AddReplyToSuffix(asyncapi.CreateRequestChannelPath)
	reqMsg.Set(payload)

	// Send request
	respMsg, err := b.ctrl.RequestToCreateOperation(ctx, reqMsg)
	if err != nil {
		return 0, err
	}

	// Check error
	if respMsg.Payload.Error != nil {
		return 0, fmt.Errorf("%d Error: %s", respMsg.Payload.Error.Code, respMsg.Payload.Error.Message)
	}

	return uint(respMsg.Payload.Id), nil
}

func (b Client) Subscribe(ctx context.Context, backtestID uint, exchange, pair string) error {
	// Set message
	reqMsg := asyncapi.NewSubscribeRequestMessage()
	reqMsg.Headers.ReplyTo = helpers.AddReplyToSuffix(asyncapi.SubscribeRequestChannelPath)
	reqMsg.Set(backtestID, exchange, pair)

	// Send request
	respMsg, err := b.ctrl.RequestToSubscribeOperation(ctx, reqMsg)
	if err != nil {
		return err
	}

	// Check error
	if respMsg.Payload.Error != nil {
		return fmt.Errorf("%d Error: %s", respMsg.Payload.Error.Code, respMsg.Payload.Error.Message)
	}

	return nil
}

func (b Client) Advance(ctx context.Context, backtestID uint) error {
	// Set message
	reqMsg := asyncapi.NewAdvanceRequestMessage()
	reqMsg.Headers.ReplyTo = helpers.AddReplyToSuffix(asyncapi.AdvanceRequestChannelPath)
	reqMsg.Set(backtestID)

	// Send request
	respMsg, err := b.ctrl.RequestToAdvanceOperation(ctx, reqMsg)
	if err != nil {
		return err
	}

	// Check error
	if respMsg.Payload.Error != nil {
		return fmt.Errorf("%d Error: %s", respMsg.Payload.Error.Code, respMsg.Payload.Error.Message)
	}

	return nil
}

func (b Client) CreateOrder(ctx context.Context, payload client.OrderCreationPayload) error {
	// Set message
	reqMsg := asyncapi.NewOrdersCreateRequestMessage()
	reqMsg.Headers.ReplyTo = helpers.AddReplyToSuffix(asyncapi.OrdersCreateRequestChannelPath)
	reqMsg.Set(payload)

	// Send request
	respMsg, err := b.ctrl.RequestToOrdersCreateOperation(ctx, reqMsg)
	if err != nil {
		return err
	}

	// Check error
	if respMsg.Payload.Error != nil {
		return fmt.Errorf("%d Error: %s", respMsg.Payload.Error.Code, respMsg.Payload.Error.Message)
	}

	return nil
}

func (b Client) GetAccounts(ctx context.Context, backtestID uint) (map[string]account.Account, error) {
	// Set message
	reqMsg := asyncapi.NewAccountsListRequestMessage()
	reqMsg.Headers.ReplyTo = helpers.AddReplyToSuffix(asyncapi.AccountsListRequestChannelPath)
	reqMsg.Set(backtestID)

	// Send request
	respMsg, err := b.ctrl.RequestToAccountsListOperation(ctx, reqMsg)
	if err != nil {
		return nil, err
	}

	// Check error
	if respMsg.Payload.Error != nil {
		return nil, fmt.Errorf("%d Error: %s", respMsg.Payload.Error.Code, respMsg.Payload.Error.Message)
	}

	return respMsg.ToModel(), nil
}

func (b Client) ServiceInfo(ctx context.Context) (clientPkg.ServiceInfo, error) {
	// Set message
	reqMsg := asyncapi.NewServiceInfoRequestMessage()
	reqMsg.Headers.ReplyTo = helpers.AddReplyToSuffix(asyncapi.ServiceInfoRequestChannelPath)

	// Send request
	respMsg, err := b.ctrl.RequestToServiceInfoOperation(ctx, reqMsg)
	if err != nil {
		return clientPkg.ServiceInfo{}, err
	}

	return respMsg.ToModel(), nil
}

func (b Client) Close(ctx context.Context) {
	b.ctrl.Close(ctx)
	b.broker.Close()
}
