package raw

import (
	"context"

	api "github.com/lerenn/cryptellation/v1/api/worker/go"
	temporalclient "go.temporal.io/sdk/client"
)

// ListCandlesticks calls the candlesticks list workflow.
func (c raw) ListCandlesticks(
	ctx context.Context,
	params api.ListCandlesticksWorkflowParams,
) (res api.ListCandlesticksWorkflowResults, err error) {
	workflowOptions := temporalclient.StartWorkflowOptions{
		TaskQueue: api.WorkerTaskQueueName,
	}

	// Execute workflow
	exec, err := c.temporal.ExecuteWorkflow(ctx, workflowOptions, api.ListCandlesticksWorkflowName, params)
	if err != nil {
		return api.ListCandlesticksWorkflowResults{}, err
	}

	// Get result and return
	err = exec.Get(ctx, &res)
	return res, err
}
