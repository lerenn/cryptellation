package client

import (
	"context"

	"github.com/lerenn/cryptellation/v1/api"
	temporalclient "go.temporal.io/sdk/client"
)

// Client is the client interface.
type Client interface {
	Temporal() temporalclient.Client
	Close(ctx context.Context)

	// Candlesticks
	ListCandlesticks(ctx context.Context, params api.ListCandlesticksWorkflowParams) (res api.ListCandlesticksWorkflowResults, err error)

	// Exchanges
	GetExchange(ctx context.Context, params api.GetExchangeWorkflowParams) (res api.GetExchangeWorkflowResults, err error)
	ListExchanges(ctx context.Context, params api.ListExchangesWorkflowParams) (res api.ListExchangesWorkflowResults, err error)

	// Service
	Info(ctx context.Context) (api.ServiceInfoResults, error)

	// Ticks
	ListenToTicks(
		ctx context.Context,
		registerParams api.RegisterForTicksListeningWorkflowParams,
	) (res api.RegisterForTicksListeningWorkflowResults, err error)
}
