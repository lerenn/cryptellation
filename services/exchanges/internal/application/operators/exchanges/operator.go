package exchanges

import (
	"context"
	"time"

	"github.com/digital-feather/cryptellation/services/exchanges/pkg/models/exchange"
)

type Operator interface {
	GetCached(ctx context.Context, expiration *time.Duration, names ...string) ([]exchange.Exchange, error)
}
