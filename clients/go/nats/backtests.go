package nats

import (
	"context"
	"fmt"
	"time"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/loggers"
	client "github.com/lerenn/cryptellation/clients/go"
	"github.com/lerenn/cryptellation/internal/ctrl/backtests/events"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/models/account"
	"github.com/lerenn/cryptellation/pkg/models/event"
	"github.com/lerenn/cryptellation/pkg/models/tick"
)

type Backtests struct {
	broker *nats.Controller
	ctrl   *events.UserController
	logger extensions.Logger
}

func NewBacktests(c config.NATS) (client.Backtests, error) {
	// Create a NATS Controller
	broker := nats.NewController(c.URL())

	// Create a logger
	logger := loggers.NewECS()

	// Create a new user controller
	ctrl, err := events.NewUserController(broker, events.WithLogger(logger))
	if err != nil {
		return nil, err
	}

	return Backtests{
		broker: broker,
		ctrl:   ctrl,
		logger: logger,
	}, nil
}

func (b Backtests) ListenEvents(ctx context.Context, backtestID uint) (<-chan event.Event, error) {
	ch := make(chan event.Event, 256)

	// Create callback when a tick appears
	callback := func(ctx context.Context, msg events.BacktestsEventMessage) {
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
	return ch, b.ctrl.SubscribeCryptellationBacktestsEventsID(ctx, events.CryptellationBacktestsEventsIDParameters{
		ID: int64(backtestID),
	}, callback)
}

func (b Backtests) Create(ctx context.Context, payload client.BacktestCreationPayload) (uint, error) {
	// Set message
	reqMsg := events.NewBacktestsCreateRequestMessage()
	reqMsg.Set(payload)

	// Send request
	respMsg, err := b.ctrl.WaitForCryptellationBacktestsCreateResponse(ctx, &reqMsg, func(ctx context.Context) error {
		return b.ctrl.PublishCryptellationBacktestsCreateRequest(ctx, reqMsg)
	})
	if err != nil {
		return 0, err
	}

	// Check error
	if respMsg.Payload.Error != nil {
		return 0, fmt.Errorf("%d Error: %s", respMsg.Payload.Error.Code, respMsg.Payload.Error.Message)
	}

	return uint(respMsg.Payload.ID), nil
}

func (b Backtests) Subscribe(ctx context.Context, backtestID uint, exchange, pair string) error {
	// Set message
	reqMsg := events.NewBacktestsSubscribeRequestMessage()
	reqMsg.Set(backtestID, exchange, pair)

	// Send request
	respMsg, err := b.ctrl.WaitForCryptellationBacktestsSubscribeResponse(ctx, &reqMsg, func(ctx context.Context) error {
		return b.ctrl.PublishCryptellationBacktestsSubscribeRequest(ctx, reqMsg)
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
	reqMsg := events.NewBacktestsAdvanceRequestMessage()
	reqMsg.Set(backtestID)

	// Send request
	respMsg, err := b.ctrl.WaitForCryptellationBacktestsAdvanceResponse(ctx, &reqMsg, func(ctx context.Context) error {
		return b.ctrl.PublishCryptellationBacktestsAdvanceRequest(ctx, reqMsg)
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
	reqMsg := events.NewBacktestsOrdersCreateRequestMessage()
	reqMsg.Set(payload)

	// Send request
	respMsg, err := b.ctrl.WaitForCryptellationBacktestsOrdersCreateResponse(ctx, &reqMsg, func(ctx context.Context) error {
		return b.ctrl.PublishCryptellationBacktestsOrdersCreateRequest(ctx, reqMsg)
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
	reqMsg := events.NewBacktestsAccountsListRequestMessage()
	reqMsg.Set(backtestID)

	// Send request
	respMsg, err := b.ctrl.WaitForCryptellationBacktestsAccountsListResponse(ctx, &reqMsg, func(ctx context.Context) error {
		return b.ctrl.PublishCryptellationBacktestsAccountsListRequest(ctx, reqMsg)
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

func (b Backtests) Close(ctx context.Context) {
	b.ctrl.Close(ctx)
	b.broker.Close()
}
