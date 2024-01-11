package app

import (
	"context"

	"github.com/lerenn/cryptellation/svc/exchanges/pkg/exchange"
)

type Exchanges interface {
	GetCached(ctx context.Context, names ...string) ([]exchange.Exchange, error)
}
