package client

import (
	"context"

	"github.com/lerenn/cryptellation/v1/api"
)

// Client is the client interface.
type Client interface {
	// Candlesticks
	ListCandlesticks(ctx context.Context, params api.ListCandlesticksParams) (res api.ListCandlesticksResults, err error)

	// Exchanges
	GetExchange(ctx context.Context, params api.GetExchangeParams) (res api.GetExchangeResults, err error)
	ListExchanges(ctx context.Context, params api.ListExchangesParams) (res api.ListExchangesResults, err error)

	// Service
	Info(ctx context.Context) (api.ServiceInfoResult, error)
	Close(ctx context.Context)
}
