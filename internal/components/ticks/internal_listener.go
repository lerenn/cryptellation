package ticks

import (
	"context"
	"log"
	"time"

	"github.com/lerenn/cryptellation/internal/components/ticks/ports/db"
	"github.com/lerenn/cryptellation/internal/components/ticks/ports/events"
	"github.com/lerenn/cryptellation/internal/components/ticks/ports/exchanges"
	"github.com/lerenn/cryptellation/pkg/models/tick"
)

const checkInterval = 1 * time.Second

type internalListener struct {
	DB        db.Port
	Events    events.Port
	Exchanges exchanges.Port

	ExchangeName string
	PairSymbol   string

	ticksChan     chan tick.Tick
	stopChan      chan struct{}
	nextCheckTime time.Time
}

func (l *internalListener) Run() (err error) {
	log.Printf("Starting listener for %q on %q\n", l.PairSymbol, l.ExchangeName)

	// Starting listening to symbol
	l.ticksChan, l.stopChan, err = l.Exchanges.ListenSymbol(l.ExchangeName, l.PairSymbol)
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

	// Close the Events listener when exiting
	defer l.Events.Close(context.TODO())

	for {
		t, open := <-l.ticksChan
		if t.Price != 0 && t.Price != lastPrice {
			err := l.Events.Publish(context.TODO(), t)
			if err != nil {
				log.Println("Publish error:", err)
				continue
			}
			lastPrice = t.Price
		}

		if !open {
			log.Printf("Closing %q listener on %q", l.PairSymbol, l.ExchangeName)
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
		count, err := l.DB.GetSymbolListenerSubscribers(ctx, l.ExchangeName, l.PairSymbol)
		if err != nil {
			return false, err
		}

		if count <= 0 {
			log.Println("Interrupting", l.ExchangeName, l.PairSymbol, "listener")
			l.stopChan <- struct{}{}
			return true, nil
		}

		l.nextCheckTime = time.Now().Add(checkInterval)
	}

	return false, nil
}
