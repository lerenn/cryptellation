package nats

import (
	"context"
	"fmt"
	"time"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
	asyncapi "github.com/lerenn/cryptellation/api/asyncapi/backtests"
	client "github.com/lerenn/cryptellation/clients/go"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/models/account"
	"github.com/lerenn/cryptellation/pkg/models/event"
	"github.com/lerenn/cryptellation/pkg/models/tick"
)

type Backtests struct {
	broker *nats.Controller
	ctrl   *asyncapi.UserController

	logger extensions.Logger
}

type BacktestsOption func(b *Backtests)

func NewBacktests(c config.NATS, options ...BacktestsOption) (client.Backtests, error) {
	var b Backtests

	// Execute options
	for _, option := range options {
		option(&b)
	}

	// Create a NATS Controller
	b.broker = nats.NewController(c.URL())

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

	return b, nil
}

func WithBacktestsLogger(logger extensions.Logger) BacktestsOption {
	return func(b *Backtests) {
		b.logger = logger
	}
}

func (b Backtests) ListenEvents(ctx context.Context, backtestID uint) (<-chan event.Event, error) {
	ch := make(chan event.Event, 256)

	// Create callback when a tick appears
	callback := func(ctx context.Context, msg asyncapi.BacktestsEventMessage) {
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
				Time:       time.Time(msg.Payload.Content.Time),
				Exchange:   string(msg.Payload.Content.Exchange),
				PairSymbol: string(msg.Payload.Content.PairSymbol),
				Price:      msg.Payload.Content.Price,
			}
		default:
			b.logger.Error(ctx, fmt.Sprintf("received unknown event type: %s", msg.Payload.Type))
			return
		}

		// Try to send tick or drop it
		select {
		case ch <- e:
		default:
			// Drop if it's full or closed
		}
	}

	// Listen to channel
	return ch, b.ctrl.SubscribeBacktestEvent(ctx, asyncapi.CryptellationBacktestsEventsParameters{
		Id: int64(backtestID),
	}, callback)
}

func (b Backtests) Create(ctx context.Context, payload client.BacktestCreationPayload) (uint, error) {
	// Set message
	reqMsg := asyncapi.NewCreateBacktestRequestMessage()
	reqMsg.Set(payload)

	// Send request
	respMsg, err := b.ctrl.WaitForCreateBacktestResponse(ctx, &reqMsg, func(ctx context.Context) error {
		return b.ctrl.PublishCreateBacktestRequest(ctx, reqMsg)
	})
	if err != nil {
		return 0, err
	}

	// Check error
	if respMsg.Payload.Error != nil {
		return 0, fmt.Errorf("%d Error: %s", respMsg.Payload.Error.Code, respMsg.Payload.Error.Message)
	}

	return uint(respMsg.Payload.Id), nil
}

func (b Backtests) Subscribe(ctx context.Context, backtestID uint, exchange, pair string) error {
	// Set message
	reqMsg := asyncapi.NewSubscribeBacktestRequestMessage()
	reqMsg.Set(backtestID, exchange, pair)

	// Send request
	respMsg, err := b.ctrl.WaitForSubscribeBacktestResponse(ctx, &reqMsg, func(ctx context.Context) error {
		return b.ctrl.PublishSubscribeBacktestRequest(ctx, reqMsg)
	})
	if err != nil {
		return err
	}

	// Check error
	if respMsg.Payload.Error != nil {
		return fmt.Errorf("%d Error: %s", respMsg.Payload.Error.Code, respMsg.Payload.Error.Message)
	}

	return nil
}

func (b Backtests) Advance(ctx context.Context, backtestID uint) error {
	// Set message
	reqMsg := asyncapi.NewAdvanceBacktestRequestMessage()
	reqMsg.Set(backtestID)

	// Send request
	respMsg, err := b.ctrl.WaitForAdvanceBacktestResponse(ctx, &reqMsg, func(ctx context.Context) error {
		return b.ctrl.PublishAdvanceBacktestRequest(ctx, reqMsg)
	})
	if err != nil {
		return err
	}

	// Check error
	if respMsg.Payload.Error != nil {
		return fmt.Errorf("%d Error: %s", respMsg.Payload.Error.Code, respMsg.Payload.Error.Message)
	}

	return nil
}

func (b Backtests) CreateOrder(ctx context.Context, payload client.OrderCreationPayload) error {
	// Set message
	reqMsg := asyncapi.NewCreateBacktestOrderRequestMessage()
	reqMsg.Set(payload)

	// Send request
	respMsg, err := b.ctrl.WaitForCreateBacktestOrderResponse(ctx, &reqMsg, func(ctx context.Context) error {
		return b.ctrl.PublishCreateBacktestOrderRequest(ctx, reqMsg)
	})
	if err != nil {
		return err
	}

	// Check error
	if respMsg.Payload.Error != nil {
		return fmt.Errorf("%d Error: %s", respMsg.Payload.Error.Code, respMsg.Payload.Error.Message)
	}

	return nil
}

func (b Backtests) GetAccounts(ctx context.Context, backtestID uint) (map[string]account.Account, error) {
	// Set message
	reqMsg := asyncapi.NewListBacktestAccountsRequestMessage()
	reqMsg.Set(backtestID)

	// Send request
	respMsg, err := b.ctrl.WaitForListBacktestAccountsResponse(ctx, &reqMsg, func(ctx context.Context) error {
		return b.ctrl.PublishListBacktestAccountsRequest(ctx, reqMsg)
	})
	if err != nil {
		return nil, err
	}

	// Check error
	if respMsg.Payload.Error != nil {
		return nil, fmt.Errorf("%d Error: %s", respMsg.Payload.Error.Code, respMsg.Payload.Error.Message)
	}

	return respMsg.ToModel(), nil
}

func (b Backtests) ServiceInfo(ctx context.Context) (client.ServiceInfo, error) {
	// Set message
	reqMsg := asyncapi.NewServiceInfoRequestMessage()

	// Send request
	respMsg, err := b.ctrl.WaitForServiceInfoResponse(ctx, &reqMsg, func(ctx context.Context) error {
		return b.ctrl.PublishServiceInfoRequest(ctx, reqMsg)
	})
	if err != nil {
		return client.ServiceInfo{}, err
	}

	return respMsg.ToModel(), nil
}

func (b Backtests) Close(ctx context.Context) {
	b.ctrl.Close(ctx)
	b.broker.Close()
}
