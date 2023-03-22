package app

import (
	"context"

	"github.com/digital-feather/cryptellation/pkg/types/candlestick"
)

type Controller interface {
	GetCached(ctx context.Context, payload GetCachedPayload) (*candlestick.List, error)
}
