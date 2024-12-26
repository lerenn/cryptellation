package raw

import (
	"github.com/lerenn/cryptellation/v1/api"
	"go.temporal.io/sdk/workflow"
)

// ListCandlesticks lists candlesticks from Cryptellation service.
func ListCandlesticks(
	ctx workflow.Context,
	params api.ListCandlesticksWorkflowParams,
	childWorkflowOptions *workflow.ChildWorkflowOptions,
) (result api.ListCandlesticksWorkflowResults, err error) {
	// Set options
	if childWorkflowOptions == nil {
		childWorkflowOptions = &workflow.ChildWorkflowOptions{}
	}
	ctx = workflow.WithChildOptions(ctx, *childWorkflowOptions)

	// Get candlesticks
	err = workflow.ExecuteChildWorkflow(ctx, api.ListCandlesticksWorkflowName, params).Get(ctx, &result)
	return result, err
}
