package exchanges

import (
	"context"

	"github.com/lerenn/cryptellation/pkg/models/exchange"
)

type Interface interface {
	GetCached(ctx context.Context, names ...string) ([]exchange.Exchange, error)
}
