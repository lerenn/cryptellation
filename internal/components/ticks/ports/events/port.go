// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=port.go -destination=mock.gen.go -package events

package events

import (
	"context"

	"github.com/lerenn/cryptellation/pkg/models/tick"
)

type Port interface {
	Publish(ctx context.Context, tick tick.Tick) error
	Subscribe(ctx context.Context, symbol string) (<-chan tick.Tick, error)
	Close(ctx context.Context)
}
