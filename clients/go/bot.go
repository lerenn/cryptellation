package client

import (
	"github.com/lerenn/cryptellation/pkg/event"
)

type Bot interface {
	OnInit(run *Run)
	OnEvent(event.Event) error
	OnExit() error

	TicksToListen() []event.TickSubscription
}
