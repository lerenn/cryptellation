package nats

import (
	"context"
	"fmt"
	"log"
	"time"

	asyncapi "github.com/digital-feather/cryptellation/api/asyncapi/backtests"
	client "github.com/digital-feather/cryptellation/clients/go"
	"github.com/digital-feather/cryptellation/pkg/account"
	"github.com/digital-feather/cryptellation/pkg/config"
	"github.com/digital-feather/cryptellation/pkg/event"
	"github.com/digital-feather/cryptellation/pkg/tick"
	"github.com/nats-io/nats.go"
)

type Backtests struct {
	nats *nats.Conn
	ctrl *asyncapi.ClientController
}

func NewBacktests(c config.NATS) (client.Backtests, error) {
	conn, err := nats.Connect(c.URL())
	if err != nil {
		return nil, err
	}

	ctrl, err := asyncapi.NewClientController(asyncapi.NewNATSController(conn))
	if err != nil {
		return nil, err
	}

	return Backtests{
		nats: conn,
		ctrl: ctrl,
	}, nil
}

func (b Backtests) ListenEvents(ctx context.Context, backtestID uint) (<-chan event.Event, error) {
	ch := make(chan event.Event, 256)

	// Create callback when a tick appears
	callback := func(msg asyncapi.BacktestsEventMessage, done bool) {
		// Check if done
		if done {
			close(ch)
			return
		}

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
			log.Printf("received unknown event type: %s", msg.Payload.Type)
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
	return ch, b.ctrl.SubscribeCryptellationBacktestsEventsID(asyncapi.CryptellationBacktestsEventsIDParameters{
		ID: int64(backtestID),
	}, callback)
}

func (b Backtests) Create(ctx context.Context, payload client.BacktestCreationPayload) (uint, error) {
	// Set message
	reqMsg := asyncapi.NewBacktestsCreateRequestMessage()
	reqMsg.Set(payload)

	// Send request
	respMsg, err := b.ctrl.WaitForCryptellationBacktestsCreateResponse(ctx, reqMsg, func() error {
		return b.ctrl.PublishCryptellationBacktestsCreateRequest(reqMsg)
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
	reqMsg := asyncapi.NewBacktestsSubscribeRequestMessage()
	reqMsg.Set(backtestID, exchange, pair)

	// Send request
	respMsg, err := b.ctrl.WaitForCryptellationBacktestsSubscribeResponse(ctx, reqMsg, func() error {
		return b.ctrl.PublishCryptellationBacktestsSubscribeRequest(reqMsg)
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
	reqMsg := asyncapi.NewBacktestsAdvanceRequestMessage()
	reqMsg.Set(backtestID)

	// Send request
	respMsg, err := b.ctrl.WaitForCryptellationBacktestsAdvanceResponse(ctx, reqMsg, func() error {
		return b.ctrl.PublishCryptellationBacktestsAdvanceRequest(reqMsg)
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
	reqMsg := asyncapi.NewBacktestsOrdersCreateRequestMessage()
	reqMsg.Set(payload)

	// Send request
	respMsg, err := b.ctrl.WaitForCryptellationBacktestsOrdersCreateResponse(ctx, reqMsg, func() error {
		return b.ctrl.PublishCryptellationBacktestsOrdersCreateRequest(reqMsg)
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
	reqMsg := asyncapi.NewBacktestsAccountsListRequestMessage()
	reqMsg.Set(backtestID)

	// Send request
	respMsg, err := b.ctrl.WaitForCryptellationBacktestsAccountsListResponse(ctx, reqMsg, func() error {
		return b.ctrl.PublishCryptellationBacktestsAccountsListRequest(reqMsg)
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

func (b Backtests) Close() {
	b.ctrl.Close()
	b.nats.Close()
}
