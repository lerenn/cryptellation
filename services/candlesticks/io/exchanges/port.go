// Generate code for mock
//go:generate go run github.com/golang/mock/mockgen -source=port.go -destination=mock.gen.go -package exchanges

package exchanges

import (
	"context"

	"github.com/digital-feather/cryptellation/pkg/candlestick"
)

type Port interface {
	GetCandlesticks(ctx context.Context, payload GetCandlesticksPayload) (*candlestick.List, error)
}
