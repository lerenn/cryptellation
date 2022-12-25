package ticks

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/digital-feather/cryptellation/services/ticks/internal/application/ports/pubsub"
	"github.com/digital-feather/cryptellation/services/ticks/internal/application/ports/vdb"
	"github.com/digital-feather/cryptellation/services/ticks/pkg/models/tick"
)

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
	db         vdb.Adapter
	ps         pubsub.Adapter
	ticksChan  chan tick.Tick
	stopChan   chan struct{}
}

const checkInterval = 10 * time.Second

func newListener(payload newListenerPayload) {
	ctx := context.Background()
	nextCheckTime := time.Now().Add(checkInterval)
	lastPrice := float64(0.0)

	for {
		t, closed := <-payload.ticksChan
		if closed {
			payload.ps.Close()
			return
		} else if lastPrice == t.Price {
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
