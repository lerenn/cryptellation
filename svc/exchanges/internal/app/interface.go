package app

import (
	"context"

	"cryptellation/svc/exchanges/pkg/exchange"
)

type Exchanges interface {
	GetCached(ctx context.Context, names ...string) ([]exchange.Exchange, error)
}
