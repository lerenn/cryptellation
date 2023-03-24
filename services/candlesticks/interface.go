package candlesticks

import (
	"context"

	"github.com/digital-feather/cryptellation/pkg/candlestick"
)

type Interface interface {
	GetCached(ctx context.Context, payload GetCachedPayload) (*candlestick.List, error)
}
