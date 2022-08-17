package pubsub

import (
	"github.com/digital-feather/cryptellation/services/ticks/internal/domain/tick"
)

type Port interface {
	Publish(tick tick.Tick) error
	Subscribe(symbol string) (<-chan tick.Tick, error)
	Close()
}
