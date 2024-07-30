package client

import (
	"cryptellation/pkg/models/event"

	"cryptellation/svc/ticks/pkg/tick"
)

type Bot interface {
	OnInit(run *Run)
	OnTick(tick.Tick) error
	OnExit() error

	TicksToListen() []event.TickSubscription
}
