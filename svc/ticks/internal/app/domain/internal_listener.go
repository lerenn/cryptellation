package domain

import (
	"context"
	"log"
	"time"

	"github.com/lerenn/cryptellation/svc/ticks/internal/app/ports/db"
	"github.com/lerenn/cryptellation/svc/ticks/internal/app/ports/events"
	"github.com/lerenn/cryptellation/svc/ticks/internal/app/ports/exchanges"
	"github.com/lerenn/cryptellation/svc/ticks/pkg/tick"
)

const checkInterval = 1 * time.Second

type internalListener struct {
	DB        db.Port
	Events    events.Port
	Exchanges exchanges.Port

	Exchange string
	Pair     string

	ticksChan     chan tick.Tick
	stopChan      chan struct{}
	nextCheckTime time.Time
}

func (l *internalListener) Run() (err error) {
	log.Printf("Starting listener for %q on %q\n", l.Pair, l.Exchange)

	// Starting listening to symbol
	l.ticksChan, l.stopChan, err = l.Exchanges.ListenSymbol(l.Exchange, l.Pair)
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
			log.Printf("Closing %q listener on %q", l.Pair, l.Exchange)
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
		count, err := l.DB.GetSymbolListenerSubscribers(ctx, l.Exchange, l.Pair)
		if err != nil {
			return false, err
		}

		if count <= 0 {
			log.Println("Interrupting", l.Exchange, l.Pair, "listener")
			l.stopChan <- struct{}{}
			return true, nil
		}

		l.nextCheckTime = time.Now().Add(checkInterval)
	}

	return false, nil
}
