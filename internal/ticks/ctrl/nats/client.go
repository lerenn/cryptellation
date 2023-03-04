package nats

import (
	"context"
	"fmt"
	"time"

	"github.com/digital-feather/cryptellation/internal/ticks/ctrl/nats/generated"
	"github.com/digital-feather/cryptellation/pkg/config"
	"github.com/digital-feather/cryptellation/pkg/tick"
	"github.com/nats-io/nats.go"
)

type client struct {
	nats *nats.Conn
	ctrl *generated.ClientController
}

func New(c config.NATS) (Client, error) {
	conn, err := nats.Connect(c.URL())
	if err != nil {
		return nil, err
	}

	ctrl, err := generated.NewClientController(generated.NewNATSController(conn))
	if err != nil {
		return nil, err
	}

	return client{
		nats: conn,
		ctrl: ctrl,
	}, nil
}

func (c client) Register(ctx context.Context, payload TicksFilterPayload) error {
	// Generate message
	msg := generated.NewRegisteringRequestMessage()
	msg.Payload.Exchange = generated.ExchangeNameSchema(payload.ExchangeName)
	msg.Payload.Pair = generated.PairSymbolSchema(payload.PairSymbol)

	// Send message
	resp, err := c.ctrl.WaitForTicksRegisterResponse(ctx, msg, func() error {
		return c.ctrl.PublishTicksRegisterRequest(msg)
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

func (c client) Listen(ctx context.Context, payload TicksFilterPayload) (<-chan tick.Tick, error) {
	ch := make(chan tick.Tick, 256)

	// Create params for channel path
	params := generated.TicksListenExchangePairParameters{
		Exchange: payload.ExchangeName,
		Pair:     payload.PairSymbol,
	}

	// Create callback when a tick appears
	callback := func(msg generated.TickMessage, done bool) {
		// Check if done
		if done {
			close(ch)
			return
		}

		// Transform message to tick
		t := tick.Tick{
			Time:       time.Time(msg.Payload.Time),
			PairSymbol: string(msg.Payload.PairSymbol),
			Price:      msg.Payload.Price,
			Exchange:   string(msg.Payload.Exchange),
		}

		// Try to send tick or drop it
		select {
		case ch <- t:
		default:
			// Drop if it's full or closed
		}
	}

	// Listen to channel
	return ch, c.ctrl.SubscribeTicksListenExchangePair(params, callback)
}

func (c client) Unregister(ctx context.Context, payload TicksFilterPayload) error {
	// Generate message
	msg := generated.NewRegisteringRequestMessage()
	msg.Payload.Exchange = generated.ExchangeNameSchema(payload.ExchangeName)
	msg.Payload.Pair = generated.PairSymbolSchema(payload.PairSymbol)

	// Send message
	resp, err := c.ctrl.WaitForTicksUnregisterResponse(ctx, msg, func() error {
		return c.ctrl.PublishTicksUnregisterRequest(msg)
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

func (c client) Close() {
	c.ctrl.Close()
	c.nats.Close()
}
