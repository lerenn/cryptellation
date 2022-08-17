package pubsub

import (
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/event"
)

type Port interface {
	Publish(backtestID uint, event event.Event) error
	Subscribe(backtestID uint) (<-chan event.Event, error)
	Close()
}
