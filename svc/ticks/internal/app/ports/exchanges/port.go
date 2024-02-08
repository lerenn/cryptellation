// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=port.go -destination=mock.gen.go -package exchanges

package exchanges

import (
	"context"

	"github.com/lerenn/cryptellation/svc/ticks/pkg/tick"
)

type Port interface {
	ListenSymbol(ctx context.Context, exchange, symbol string) (chan tick.Tick, chan struct{}, error)
}
