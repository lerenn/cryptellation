package client

import (
	"context"

	"github.com/lerenn/cryptellation/pkg/models/event"

	"github.com/lerenn/cryptellation/ticks/pkg/tick"
)

type Bot interface {
	OnInit(ctx context.Context, run *Run)
	OnTick(ctx context.Context, t tick.Tick) error
	OnExit(ctx context.Context) error

	TicksToListen(ctx context.Context) []event.TickSubscription
}
