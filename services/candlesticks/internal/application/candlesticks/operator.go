package candlesticks

import (
	"context"

	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/candlestick"
)

type Operator interface {
	GetCached(ctx context.Context, payload GetCachedPayload) (*candlestick.List, error)
}
