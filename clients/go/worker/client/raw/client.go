package raw

import (
	"context"

	api "github.com/lerenn/cryptellation/v1/api/worker/go"
	"github.com/lerenn/cryptellation/v1/pkg/config"
	temporalclient "go.temporal.io/sdk/client"
)

// Client is the interface for the raw client (direct access to params/result
// from temporal).
type Client interface {
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
	GetBacktest(
		ctx context.Context,
		params api.GetBacktestWorkflowParams,
	) (api.GetBacktestWorkflowResults, error)
	ListBacktests(
		ctx context.Context,
		params api.ListBacktestsWorkflowParams,
	) (api.ListBacktestsWorkflowResults, error)
	SubscribeToBacktestPrice(
		ctx context.Context,
		params api.SubscribeToBacktestPriceWorkflowParams,
	) (api.SubscribeToBacktestPriceWorkflowResults, error)

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

	// Forwardtests

	CreateForwardtest(
		ctx context.Context,
		params api.CreateForwardtestWorkflowParams,
	) (api.CreateForwardtestWorkflowResults, error)
	GetForwardtestStatus(
		ctx context.Context,
		params api.GetForwardtestStatusWorkflowParams,
	) (api.GetForwardtestStatusWorkflowResults, error)
	ListForwardtests(
		ctx context.Context,
		params api.ListForwardtestsWorkflowParams,
	) (api.ListForwardtestsWorkflowResults, error)
	CreateForwardtestOrder(
		ctx context.Context,
		params api.CreateForwardtestOrderWorkflowParams,
	) (api.CreateForwardtestOrderWorkflowResults, error)
	ListForwardtestAccounts(
		ctx context.Context,
		params api.ListForwardtestAccountsWorkflowParams,
	) (api.ListForwardtestAccountsWorkflowResults, error)

	// Indicators

	ListSMA(
		ctx context.Context,
		params api.ListSMAWorkflowParams,
	) (res api.ListSMAWorkflowResults, err error)

	// Service

	Info(ctx context.Context) (api.ServiceInfoResults, error)

	// Ticks

	ListenToTicks(
		ctx context.Context,
		registerParams api.RegisterForTicksListeningWorkflowParams,
	) (res api.RegisterForTicksListeningWorkflowResults, err error)
}

var _ Client = raw{}

type raw struct {
	temporal temporalclient.Client
}

// NewClient creates a new client to execute temporal workflows.
// temporalConfig is the optional configuration to use for the temporal client.
func NewClient(temporalConfig ...config.Temporal) (Client, error) {
	var t config.Temporal

	if len(temporalConfig) > 0 {
		t = temporalConfig[0]
	}

	// Load temporal configuration
	t = config.LoadTemporal(&t)
	if err := t.Validate(); err != nil {
		return raw{}, err
	}

	// Create temporal client
	c, err := t.CreateTemporalClient()
	if err != nil {
		return raw{}, err
	}

	return &raw{temporal: c}, nil
}

func (c raw) Temporal() temporalclient.Client {
	return c.temporal
}

// Close closes the client.
func (c raw) Close(_ context.Context) {
	c.temporal.Close()
}
