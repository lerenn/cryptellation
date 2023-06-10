package candlesticks

import (
	"context"

	"github.com/lerenn/cryptellation/pkg/models/candlestick"
)

type Interface interface {
	GetCached(ctx context.Context, payload GetCachedPayload) (*candlestick.List, error)
}
