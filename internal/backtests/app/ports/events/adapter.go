// Generate code for mock
//go:generate go run github.com/golang/mock/mockgen -source=adapter.go -destination=mock.gen.go -package events

package events

import (
	"github.com/digital-feather/cryptellation/pkg/types/event"
)

type Adapter interface {
	Publish(backtestID uint, event event.Event) error
	Subscribe(backtestID uint) (<-chan event.Event, error)
	Close()
}
