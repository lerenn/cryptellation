package direct

import (
	"context"

	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/clients/go/direct/raw"
	"github.com/lerenn/cryptellation/v1/pkg/config"
	temporalclient "go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
)

// Client is the client interface.
type Client interface {
	RawClient() raw.Client
	Temporal() temporalclient.Client
	Close(ctx context.Context)

	// Backtests

	NewBacktest(
		ctx context.Context,
		params api.CreateBacktestWorkflowParams,
	) (Backtest, error)

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
		exchange, pair string,
		callback func(ctx workflow.Context, params api.ListenToTicksCallbackWorkflowParams) error,
	) error
}

type client struct {
	raw.Client
}

// NewClient creates a new client to execute temporal workflows.
// temporalConfig is the optional configuration to use for the temporal client.
func NewClient(temporalConfig ...config.Temporal) (Client, error) {
	c, err := raw.NewClient(temporalConfig...)
	if err != nil {
		return client{}, err
	}

	return client{
		Client: c,
	}, nil
}

func (c client) RawClient() raw.Client {
	return c.Client
}