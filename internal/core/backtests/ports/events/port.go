// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=port.go -destination=mock.gen.go -package events

package events

import (
	"github.com/lerenn/cryptellation/pkg/models/event"
)

type Port interface {
	Publish(backtestID uint, event event.Event) error
	Subscribe(backtestID uint) (<-chan event.Event, error)
	Close()
}
