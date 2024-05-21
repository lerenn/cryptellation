package client

import (
	"github.com/lerenn/cryptellation/pkg/models/event"
)

type Bot interface {
	OnInit(run *Run)
	OnEvent(event.Event) error
	OnExit() error

	TicksToListen() []event.TickSubscription
}
