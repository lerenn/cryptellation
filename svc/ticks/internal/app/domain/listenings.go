package domain

import (
	"context"
	"sync"
	"time"

	"cryptellation/pkg/adapters/telemetry"
	"cryptellation/pkg/models/event"
)

type subscription struct {
	TickSubscription event.TickSubscription
	Adapters         adapters
	LastRequest      time.Time
	Cancel           context.CancelFunc
}

func (s *subscription) ListenTicks(ctx context.Context, ts event.TickSubscription) {
	telemetry.L(ctx).Infof("subscribing to %q ticks", ts)
	ticks, done, err := s.Adapters.Exchanges.ListenSymbol(ctx, ts)
	if err != nil {
		telemetry.L(ctx).Errorf("error when subscribing to %q: %s", ts, err)
		return
	}
	defer func() {
		done <- struct{}{}
	}()

	for {
		select {
		case <-ctx.Done():
			telemetry.L(ctx).Infof("context done, stop listening ticks for %q", ts)
			return
		case tick, ok := <-ticks:
			if !ok {
				telemetry.L(ctx).Errorf("ticks channel unexpectedly closed for %q", ts)
				return
			}

			if err := s.Adapters.Events.PublishTick(ctx, tick); err != nil {
				telemetry.L(ctx).Errorf("error when publishing tick on %q: %s", ts, err)
				return
			}
		}
	}
}

type listenings struct {
	adapters      adapters
	subscriptions map[event.TickSubscription]*subscription
	lock          sync.Mutex
}

func newListenings(adapters adapters) listenings {
	return listenings{
		adapters:      adapters,
		subscriptions: make(map[event.TickSubscription]*subscription),
	}
}

func (l *listenings) UpdateLastNotificationSeen(ts event.TickSubscription) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if sub, ok := l.subscriptions[ts]; ok {
		telemetry.L(context.Background()).Infof("updating last notification seen for %q", ts)
		sub.LastRequest = time.Now()
		return
	}

	cancelableCtx, cancel := context.WithCancel(context.Background())
	sub := &subscription{
		Adapters:         l.adapters,
		TickSubscription: ts,
		LastRequest:      time.Now(),
		Cancel:           cancel,
	}
	go sub.ListenTicks(cancelableCtx, ts)
	go l.watchNoListener(cancelableCtx, sub)

	l.subscriptions[ts] = sub
}

func (l *listenings) watchNoListener(ctx context.Context, sub *subscription) {
	// Remove subscription at the end of the function
	defer l.removeSubscription(sub.TickSubscription)

	// Wait for cancellation or a timeout
	for {
		select {
		case <-ctx.Done():
			telemetry.L(ctx).Infof("context canceled for watching listener for %q", sub.TickSubscription)
			return
		case <-time.After(10 * time.Second):
			if time.Since(sub.LastRequest) < 10*time.Second {
				continue
			}

			telemetry.L(ctx).Infof("no request for %q in the last 10 seconds", sub.TickSubscription)
			return
		}
	}
}

func (l *listenings) removeSubscription(ts event.TickSubscription) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if sub, ok := l.subscriptions[ts]; ok {
		telemetry.L(context.Background()).Infof("canceling and deleting subscription for %q", ts)
		sub.Cancel()
		delete(l.subscriptions, ts)
	}
}

func (l *listenings) Close(ctx context.Context) {
	for _, sub := range l.subscriptions {
		l.removeSubscription(sub.TickSubscription)
	}
}
