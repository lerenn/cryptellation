package ticks

import (
	"context"
	"log"
	"time"

	"github.com/digital-feather/cryptellation/services/ticks/internal/application/ports/exchanges"
	"github.com/digital-feather/cryptellation/services/ticks/internal/application/ports/pubsub"
	"github.com/digital-feather/cryptellation/services/ticks/internal/application/ports/vdb"
	"github.com/digital-feather/cryptellation/services/ticks/pkg/models/tick"
)

const checkInterval = 10 * time.Second

type internalListener struct {
	DB         vdb.Adapter
	PubSub     pubsub.Adapter
	Exchange   exchanges.Adapter
	PairSymbol string

	ticksChan     chan tick.Tick
	stopChan      chan struct{}
	nextCheckTime time.Time
}

func (l *internalListener) Run() (err error) {
	// Starting listening to symbol
	l.ticksChan, l.stopChan, err = l.Exchange.ListenSymbol(l.PairSymbol)
	if err != nil {
		return err
	}

	// Setting next check time for loop
	l.nextCheckTime = time.Now().Add(checkInterval)

	// Launching internal loop for listening
	go l.internalLoop()

	return nil
}

func (l *internalListener) internalLoop() {
	lastPrice := float64(0.0)

	// Close the pubsub listener when exiting
	defer l.PubSub.Close()

	for {
		t, closed := <-l.ticksChan
		if t.Price != 0 || t.Price != lastPrice {
			err := l.PubSub.Publish(t)
			if err != nil {
				log.Println("Publish error:", err)
				continue
			}
			lastPrice = t.Price
		}

		if closed {
			break
		}

		if finished, err := l.setNextCheckTimeIfNeeded(); err != nil {
			log.Println(err)
			continue
		} else if finished {
			break
		}
	}
}

func (l *internalListener) setNextCheckTimeIfNeeded() (finished bool, err error) {
	ctx := context.Background()

	if l.nextCheckTime.Before(time.Now()) {
		count, err := l.DB.GetSymbolListenerCount(ctx, l.Exchange.Name(), l.PairSymbol)
		if err != nil {
			return false, err
		}

		if count <= 0 {
			log.Println("Interrupting", l.Exchange.Name(), l.PairSymbol, "listener")
			l.stopChan <- struct{}{}
			return true, nil
		}

		l.nextCheckTime = time.Now().Add(checkInterval)
	}

	return false, nil
}
