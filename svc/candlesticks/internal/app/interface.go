package app

import (
	"context"
	"time"

	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"
)

type Candlesticks interface {
	GetCached(ctx context.Context, payload GetCachedPayload) (*candlestick.List, error)
}

type GetCachedPayload struct {
	Exchange string
	Pair     string
	Period   period.Symbol
	Start    *time.Time
	End      *time.Time
	Limit    uint
}
