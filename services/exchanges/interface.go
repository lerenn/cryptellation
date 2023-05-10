package exchanges

import (
	"context"

	"github.com/lerenn/cryptellation/pkg/exchange"
)

type Interface interface {
	GetCached(ctx context.Context, names ...string) ([]exchange.Exchange, error)
}
