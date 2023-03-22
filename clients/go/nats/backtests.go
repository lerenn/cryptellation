package nats

import (
	"context"
	"log"
	"time"

	client "github.com/digital-feather/cryptellation/clients/go"
	"github.com/digital-feather/cryptellation/internal/backtests/infra/events/nats/generated"
	"github.com/digital-feather/cryptellation/pkg/config"
	"github.com/digital-feather/cryptellation/pkg/types/event"
	"github.com/digital-feather/cryptellation/pkg/types/tick"
	"github.com/nats-io/nats.go"
)

type Backtests struct {
	nats *nats.Conn
	ctrl *generated.ClientController
}

func NewBacktests(c config.NATS) (client.Backtests, error) {
	conn, err := nats.Connect(c.URL())
	if err != nil {
		return nil, err
	}

	ctrl, err := generated.NewClientController(generated.NewNATSController(conn))
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
	callback := func(msg generated.BacktestsEventMessage, done bool) {
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
	return ch, b.ctrl.SubscribeBacktestsEventsID(generated.BacktestsEventsIDParameters{
		ID: int64(backtestID),
	}, callback)
}

func (b Backtests) Create(ctx context.Context, payload client.BacktestCreationPayload) (int, error) {
	// Set message
	msg := generated.NewBacktestsCreateRequestMessage()
	msg.Payload.StartTime = generated.DateSchema(payload.StartTime)
	msg.Payload.EndTime = (*generated.DateSchema)(payload.EndTime)
	// TODO

	return 0, nil
}

// func accountModelsToAPI(accounts map[string]account.Account) []generated.AccountSchema {
// 	apiAccounts := make([]generated.AccountSchema, 0, len(accounts))
// 	for name, acc := range accounts {
// 		apiAccounts = append(apiAccounts, generated.AccountSchema{
// 			Name: name,
// 		})
// 	}
// }

func (b Backtests) Close() {
	b.ctrl.Close()
	b.nats.Close()
}
