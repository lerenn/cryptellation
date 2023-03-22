package app

import (
	"context"

	"github.com/digital-feather/cryptellation/pkg/types/exchange"
)

type Controller interface {
	GetCached(ctx context.Context, names ...string) ([]exchange.Exchange, error)
}
