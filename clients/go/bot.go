package client

import (
	"github.com/lerenn/cryptellation/pkg/models/event"
	"github.com/lerenn/cryptellation/svc/ticks/pkg/tick"
)

type Bot interface {
	OnInit(run *Run)
	OnTick(tick.Tick) error
	OnExit() error

	TicksToListen() []event.TickSubscription
}
