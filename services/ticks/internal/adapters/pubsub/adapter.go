// Generate code for mock
//go:generate go run -mod=mod github.com/golang/mock/mockgen -source=adapter.go -destination=mock.gen.go -package pubsub

package pubsub

import (
	"github.com/digital-feather/cryptellation/services/ticks/pkg/models/tick"
)

type Adapter interface {
	Publish(tick tick.Tick) error
	Subscribe(symbol string) (<-chan tick.Tick, error)
	Close()
}
