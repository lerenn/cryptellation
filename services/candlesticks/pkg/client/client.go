package client

import (
	"context"

	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/candlestick"
)

type Client interface {
	ReadCandlesticks(ctx context.Context, payload ReadCandlestickPayload) (*candlestick.List, error)
}
