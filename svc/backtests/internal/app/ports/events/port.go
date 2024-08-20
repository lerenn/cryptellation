// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=port.go -destination=mock.gen.go -package events

package events

import (
	"context"

	"github.com/lerenn/cryptellation/pkg/models/event"

	"github.com/google/uuid"
)

type Port interface {
	Publish(ctx context.Context, backtestID uuid.UUID, event event.Event) error
	Subscribe(ctx context.Context, backtestID uuid.UUID) (<-chan event.Event, error)
	Close(ctx context.Context)
}
