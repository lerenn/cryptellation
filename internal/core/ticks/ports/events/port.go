// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=port.go -destination=mock.gen.go -package events

package events

import (
	"github.com/lerenn/cryptellation/pkg/models/tick"
)

type Port interface {
	Publish(tick tick.Tick) error
	Subscribe(symbol string) (<-chan tick.Tick, error)
	Close()
}
