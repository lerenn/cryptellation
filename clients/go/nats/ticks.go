package nats

import (
	"context"
	"fmt"
	"time"

	client "github.com/digital-feather/cryptellation/clients/go"
	"github.com/digital-feather/cryptellation/internal/ticks/infra/events/nats/generated"
	"github.com/digital-feather/cryptellation/pkg/config"
	"github.com/digital-feather/cryptellation/pkg/types/tick"
	"github.com/nats-io/nats.go"
)

type Ticks struct {
	nats *nats.Conn
	ctrl *generated.ClientController
}

func NewTicks(c config.NATS) (client.Ticks, error) {
	conn, err := nats.Connect(c.URL())
	if err != nil {
		return nil, err
	}

	ctrl, err := generated.NewClientController(generated.NewNATSController(conn))
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
	msg := generated.NewRegisteringRequestMessage()
	msg.Payload.Exchange = generated.ExchangeNameSchema(payload.ExchangeName)
	msg.Payload.Pair = generated.PairSymbolSchema(payload.PairSymbol)

	// Send message
	resp, err := t.ctrl.WaitForTicksRegisterResponse(ctx, msg, func() error {
		return t.ctrl.PublishTicksRegisterRequest(msg)
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
	params := generated.TicksListenExchangePairParameters{
		Exchange: generated.ExchangeNameSchema(payload.ExchangeName),
		Pair:     generated.PairSymbolSchema(payload.PairSymbol),
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
	return ch, t.ctrl.SubscribeTicksListenExchangePair(params, callback)
}

func (t Ticks) Unregister(ctx context.Context, payload client.TicksFilterPayload) error {
	// Generate message
	msg := generated.NewRegisteringRequestMessage()
	msg.Payload.Exchange = generated.ExchangeNameSchema(payload.ExchangeName)
	msg.Payload.Pair = generated.PairSymbolSchema(payload.PairSymbol)

	// Send message
	resp, err := t.ctrl.WaitForTicksUnregisterResponse(ctx, msg, func() error {
		return t.ctrl.PublishTicksUnregisterRequest(msg)
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
