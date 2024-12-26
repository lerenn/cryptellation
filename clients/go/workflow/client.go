package wfclient

import (
	"github.com/lerenn/cryptellation/v1/api"
	"go.temporal.io/sdk/workflow"
)

// Client is the Cryptellation client to use inside workflows.
type Client interface {
	// Candlesticks

	ListCandlesticks(
		ctx workflow.Context,
		params api.ListCandlesticksWorkflowParams,
		childWorkflowOptions *workflow.ChildWorkflowOptions,
	) (result api.ListCandlesticksWorkflowResults, err error)

	// Exchanges

	GetExchange(
		ctx workflow.Context,
		params api.GetExchangeWorkflowParams,
		childWorkflowOptions *workflow.ChildWorkflowOptions,
	) (result api.GetExchangeWorkflowResults, err error)

	// Run (backtests, forwardtests, live)

	SubscribeToPrice(ctx workflow.Context, params SubscribeToPriceParams) error
}

type client struct{}

// NewClient creates a new client to execute temporal workflows.
// temporalConfig is the optional configuration to use for the temporal client.
func NewClient() Client {
	return client{}
}
