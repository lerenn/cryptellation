package pubsub

import (
	"github.com/digital-feather/cryptellation/services/ticks/pkg/models/tick"
)

type Adapter interface {
	Publish(tick tick.Tick) error
	Subscribe(symbol string) (<-chan tick.Tick, error)
	Close()
}
