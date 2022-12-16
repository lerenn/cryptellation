// Generate code for mock
//go:generate go run -mod=mod github.com/golang/mock/mockgen -source=adapter.go -destination=mock.gen.go -package pubsub

package pubsub

import (
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/event"
)

type Adapter interface {
	Publish(backtestID uint, event event.Event) error
	Subscribe(backtestID uint) (<-chan event.Event, error)
	Close()
}
