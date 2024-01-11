// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=port.go -destination=mock.gen.go -package exchanges

package exchanges

import (
	"context"

	"github.com/lerenn/cryptellation/pkg/adapters/exchanges/port"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
)

type Port interface {
	GetCandlesticks(ctx context.Context, payload port.GetCandlesticksPayload) (*candlestick.List, error)
}
