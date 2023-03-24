package exchanges

import (
	"context"

	"github.com/digital-feather/cryptellation/pkg/exchange"
)

type Interface interface {
	GetCached(ctx context.Context, names ...string) ([]exchange.Exchange, error)
}
