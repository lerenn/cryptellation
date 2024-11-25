package client

import (
	"context"

	"github.com/lerenn/cryptellation/v1/api"
	temporalclient "go.temporal.io/sdk/client"
)

// ListCandlesticks calls the candlesticks list workflow.
func (c client) ListCandlesticks(
	ctx context.Context,
	params api.ListCandlesticksParams,
) (res api.ListCandlesticksResults, err error) {
	workflowOptions := temporalclient.StartWorkflowOptions{
		TaskQueue: api.WorkerTaskQueueName,
	}

	// Execute workflow
	exec, err := c.temporal.ExecuteWorkflow(ctx, workflowOptions, api.ListCandlesticksWorkflowName, params)
	if err != nil {
		return api.ListCandlesticksResults{}, err
	}

	// Get result and return
	err = exec.Get(ctx, &res)
	return res, err
}
