package client

import (
	"context"

	"github.com/lerenn/cryptellation/v1/api"
)

// Client is the client interface.
type Client interface {
	// Candlesticks
	ListCandlesticks(ctx context.Context, params api.ListCandlesticksWorkflowParams) (res api.ListCandlesticksWorkflowResults, err error)

	// Exchanges
	GetExchange(ctx context.Context, params api.GetExchangeWorkflowParams) (res api.GetExchangeWorkflowResults, err error)
	ListExchanges(ctx context.Context, params api.ListExchangesWorkflowParams) (res api.ListExchangesWorkflowResults, err error)

	// Service
	Info(ctx context.Context) (api.ServiceInfoResults, error)
	Close(ctx context.Context)

	// Ticks
	ListenToTicks(ctx context.Context, params api.RegisterForTicksListeningWorkflowParams) (res api.RegisterForTicksListeningWorkflowResults, err error)
}
