package app

import (
	"context"

	"github.com/lerenn/cryptellation/exchanges/pkg/exchange"
)

type Exchanges interface {
	GetCached(ctx context.Context, names ...string) ([]exchange.Exchange, error)
}
