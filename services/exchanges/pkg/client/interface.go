package client

import (
	"context"

	"github.com/digital-feather/cryptellation/services/exchanges/pkg/models/exchange"
)

type Client interface {
	ReadExchanges(ctx context.Context, names ...string) ([]exchange.Exchange, error)
}
