package client

import (
	"context"
	"cryptellation/pkg/models/event"

	"cryptellation/svc/ticks/pkg/tick"
)

type Bot interface {
	OnInit(ctx context.Context, run *Run)
	OnTick(ctx context.Context, t tick.Tick) error
	OnExit(ctx context.Context) error

	TicksToListen(ctx context.Context) []event.TickSubscription
}
