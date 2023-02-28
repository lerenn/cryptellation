package app

import (
	"context"

	"github.com/digital-feather/cryptellation/pkg/candlestick"
)

type Controller interface {
	GetCached(ctx context.Context, payload GetCachedPayload) (*candlestick.List, error)
}
