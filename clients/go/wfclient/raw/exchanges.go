package raw

import (
	"github.com/lerenn/cryptellation/v1/api"
	"go.temporal.io/sdk/workflow"
)

// GetExchange gets exchange info from Cryptellation service.
func GetExchange(
	ctx workflow.Context,
	params api.GetExchangeWorkflowParams,
	childWorkflowOptions *workflow.ChildWorkflowOptions,
) (result api.GetExchangeWorkflowResults, err error) {
	// Set options
	if childWorkflowOptions == nil {
		childWorkflowOptions = &workflow.ChildWorkflowOptions{}
	}
	ctx = workflow.WithChildOptions(ctx, *childWorkflowOptions)

	// Get exchange info
	err = workflow.ExecuteChildWorkflow(ctx, api.GetExchangeWorkflowName, params).Get(ctx, &result)
	return result, err
}
