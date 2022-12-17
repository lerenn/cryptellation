// Generate code for mock
//go:generate go run -mod=mod github.com/golang/mock/mockgen -source=adapter.go -destination=mock.gen.go -package exchanges

package exchanges

import (
	"context"

	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/candlestick"
)

type Adapter interface {
	GetCandlesticks(ctx context.Context, payload GetCandlesticksPayload) (*candlestick.List, error)
}
