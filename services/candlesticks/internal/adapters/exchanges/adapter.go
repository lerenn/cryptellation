package exchanges

import (
	"context"

	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/candlestick"
)

type Adapter interface {
	GetCandlesticks(ctx context.Context, payload GetCandlesticksPayload) (*candlestick.List, error)
}
