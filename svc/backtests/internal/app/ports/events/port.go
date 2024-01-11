// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=port.go -destination=mock.gen.go -package events

package events

import (
	"context"

	"github.com/lerenn/cryptellation/svc/backtests/pkg/event"
)

type Port interface {
	Publish(ctx context.Context, backtestID uint, event event.Event) error
	Subscribe(ctx context.Context, backtestID uint) (<-chan event.Event, error)
	Close(ctx context.Context)
}
