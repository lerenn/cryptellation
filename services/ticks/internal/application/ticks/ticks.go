package ticks

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/digital-feather/cryptellation/services/ticks/internal/adapters/exchanges"
	"github.com/digital-feather/cryptellation/services/ticks/internal/adapters/pubsub"
	"github.com/digital-feather/cryptellation/services/ticks/internal/adapters/vdb"
	"github.com/digital-feather/cryptellation/services/ticks/pkg/models/tick"
)

type Ticks struct {
	vdb       vdb.Port
	pubsub    pubsub.Port
	exchanges map[string]exchanges.Port
}

func New(ps pubsub.Port, db vdb.Port, exchanges map[string]exchanges.Port) *Ticks {
	if ps == nil {
		panic("nil pubsub")
	}

	if db == nil {
		panic("nil vdb")
	}

	if exchanges == nil {
		panic("nil exchanges clients")
	}

	return &Ticks{
		pubsub:    ps,
		exchanges: exchanges,
		vdb:       db,
	}
}

func (t Ticks) Listen(exchange, pairSymbol string) (<-chan tick.Tick, error) {
	return t.pubsub.Subscribe(pairSymbol)
}

func (t Ticks) Register(ctx context.Context, exchange, pairSymbol string) (int64, error) {
	count, err := t.vdb.IncrementSymbolListenerCount(ctx, exchange, pairSymbol)
	if err != nil {
		return count, err
	}

	if count == 1 {
		err := t.launchListener(exchange, pairSymbol)
		if err != nil {
			return count, err
		}
	}

	return count, nil
}

func (t Ticks) launchListener(exchange, pairSymbol string) error {
	exch, exists := t.exchanges[exchange]
	if !exists {
		return fmt.Errorf("exchange %q doesn't exists", exchange)
	}

	ticksChan, stopChan, err := exch.ListenSymbol(pairSymbol)
	if err != nil {
		return fmt.Errorf("listening to symbol: %w", err)
	}

	go newListener(newListenerPayload{
		exchange:   exchange,
		pairSymbol: pairSymbol,
		db:         t.vdb,
		ps:         t.pubsub,
		ticksChan:  ticksChan,
		stopChan:   stopChan,
	})

	return nil
}

type newListenerPayload struct {
	exchange   string
	pairSymbol string
	db         vdb.Port
	ps         pubsub.Port
	ticksChan  chan tick.Tick
	stopChan   chan struct{}
}

const checkInterval = 10 * time.Second

func newListener(payload newListenerPayload) {
	ctx := context.Background()
	nextCheckTime := time.Now().Add(checkInterval)
	lastPrice := float64(0.0)

	for {
		t := <-payload.ticksChan
		if lastPrice == t.Price {
			continue
		}
		lastPrice = t.Price

		err := payload.ps.Publish(t)
		if err != nil {
			log.Println("Publish error:", err)
			continue
		}

		if nextCheckTime.Before(time.Now()) {
			count, err := payload.db.GetSymbolListenerCount(ctx, payload.exchange, payload.pairSymbol)
			if err != nil {
				log.Println(err)
				continue
			}

			if count <= 0 {
				log.Println("Interrupting", payload.exchange, payload.pairSymbol, "listener")
				payload.stopChan <- struct{}{}
				break
			}

			nextCheckTime = time.Now().Add(checkInterval)
		}
	}
}

func (t Ticks) Unregister(ctx context.Context, exchange, pairSymbol string) (int64, error) {
	count, err := t.vdb.DecrementSymbolListenerCount(ctx, exchange, pairSymbol)
	if err != nil {
		return count, err
	}

	// TODO Unregister listener

	return count, nil
}
