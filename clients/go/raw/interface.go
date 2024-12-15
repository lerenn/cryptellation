package raw

import (
	"context"

	"github.com/lerenn/cryptellation/v1/api"
	temporalclient "go.temporal.io/sdk/client"
)

// Raw is the interface for the raw client (direct access to params/result
// from temporal).
type Raw interface {
	Temporal() temporalclient.Client
	Close(ctx context.Context)

	// Backtests

	CreateBacktest(
		ctx context.Context,
		params api.CreateBacktestWorkflowParams,
	) (api.CreateBacktestWorkflowResults, error)
	RunBacktest(
		ctx context.Context,
		params api.RunBacktestWorkflowParams,
	) (api.RunBacktestWorkflowResults, error)

	// Candlesticks

	ListCandlesticks(
		ctx context.Context,
		params api.ListCandlesticksWorkflowParams,
	) (res api.ListCandlesticksWorkflowResults, err error)

	// Exchanges

	GetExchange(
		ctx context.Context,
		params api.GetExchangeWorkflowParams,
	) (res api.GetExchangeWorkflowResults, err error)
	ListExchanges(
		ctx context.Context,
		params api.ListExchangesWorkflowParams,
	) (res api.ListExchangesWorkflowResults, err error)

	// Service

	Info(ctx context.Context) (api.ServiceInfoResults, error)

	// Ticks

	ListenToTicks(
		ctx context.Context,
		registerParams api.RegisterForTicksListeningWorkflowParams,
	) (res api.RegisterForTicksListeningWorkflowResults, err error)
}
