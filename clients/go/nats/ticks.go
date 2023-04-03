package nats

import (
	"context"
	"fmt"

	asyncapi "github.com/digital-feather/cryptellation/api/asyncapi/ticks"
	client "github.com/digital-feather/cryptellation/clients/go"
	"github.com/digital-feather/cryptellation/pkg/config"
	"github.com/digital-feather/cryptellation/pkg/tick"
	"github.com/nats-io/nats.go"
)

type Ticks struct {
	nats *nats.Conn
	ctrl *asyncapi.ClientController
}

func NewTicks(c config.NATS) (client.Ticks, error) {
	conn, err := nats.Connect(c.URL())
	if err != nil {
		return nil, err
	}

	ctrl, err := asyncapi.NewClientController(asyncapi.NewNATSController(conn))
	if err != nil {
		return nil, err
	}

	return Ticks{
		nats: conn,
		ctrl: ctrl,
	}, nil
}

func (t Ticks) Register(ctx context.Context, payload client.TicksFilterPayload) error {
	// Generate message
	msg := asyncapi.NewRegisteringRequestMessage()
	msg.Set(payload)

	// Send message
	resp, err := t.ctrl.WaitForCryptellationTicksRegisterResponse(ctx, msg, func() error {
		return t.ctrl.PublishCryptellationTicksRegisterRequest(msg)
	})
	if err != nil {
		return err
	}

	// Check error from server
	if resp.Payload.Error != nil {
		return fmt.Errorf("%d Error: %s", resp.Payload.Error.Code, resp.Payload.Error.Message)
	}

	return nil
}

func (t Ticks) Listen(ctx context.Context, payload client.TicksFilterPayload) (<-chan tick.Tick, error) {
	ch := make(chan tick.Tick, 256)

	// Create params for channel path
	params := asyncapi.CryptellationTicksListenExchangePairParameters{
		Exchange: asyncapi.ExchangeNameSchema(payload.ExchangeName),
		Pair:     asyncapi.PairSymbolSchema(payload.PairSymbol),
	}

	// Create callback when a tick appears
	callback := func(msg asyncapi.TickMessage, done bool) {
		// Check if done
		if done {
			close(ch)
			return
		}

		// Try to send tick or drop it
		select {
		case ch <- msg.ToModel():
		default:
			// Drop if it's full or closed
		}
	}

	// Listen to channel
	return ch, t.ctrl.SubscribeCryptellationTicksListenExchangePair(params, callback)
}

func (t Ticks) Unregister(ctx context.Context, payload client.TicksFilterPayload) error {
	// Generate message
	msg := asyncapi.NewRegisteringRequestMessage()
	msg.Set(payload)

	// Send message
	resp, err := t.ctrl.WaitForCryptellationTicksUnregisterResponse(ctx, msg, func() error {
		return t.ctrl.PublishCryptellationTicksUnregisterRequest(msg)
	})
	if err != nil {
		return err
	}

	// Check error from server
	if resp.Payload.Error != nil {
		return fmt.Errorf("%d Error: %s", resp.Payload.Error.Code, resp.Payload.Error.Message)
	}

	return nil
}

func (t Ticks) Close() {
	t.ctrl.Close()
	t.nats.Close()
}
