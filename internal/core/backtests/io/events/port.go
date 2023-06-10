// Generate code for mock
//go:generate go run github.com/golang/mock/mockgen -source=port.go -destination=mock.gen.go -package events

package events

import (
	"github.com/lerenn/cryptellation/pkg/models/event"
)

type Port interface {
	Publish(backtestID uint, event event.Event) error
	Subscribe(backtestID uint) (<-chan event.Event, error)
	Close()
}
