package raw

import (
	api "github.com/lerenn/cryptellation/v1/api/worker/go"
	"go.temporal.io/sdk/workflow"
)

// SubscribeToBacktestPrice subscribes to the backtest price.
func SubscribeToBacktestPrice(
	ctx workflow.Context,
	params api.SubscribeToBacktestPriceWorkflowParams,
) (api.SubscribeToBacktestPriceWorkflowResults, error) {
	// Set options
	childWorkflowOptions := workflow.ChildWorkflowOptions{}
	ctx = workflow.WithChildOptions(ctx, childWorkflowOptions)

	// Execute child workflow
	var res api.SubscribeToBacktestPriceWorkflowResults
	err := workflow.ExecuteChildWorkflow(ctx, api.SubscribeToBacktestPriceWorkflowName, params).Get(ctx, &res)
	if err != nil {
		return api.SubscribeToBacktestPriceWorkflowResults{}, err
	}

	return res, nil
}
