package client

import (
	"context"

	"github.com/lerenn/cryptellation/v1/api"
	temporalclient "go.temporal.io/sdk/client"
)

// ListExchanges calls the exchanges list workflow.
func (c client) ListExchanges(
	ctx context.Context,
	params api.ListExchangesParams,
) (res api.ListExchangesResults, err error) {
	workflowOptions := temporalclient.StartWorkflowOptions{
		TaskQueue: api.WorkerTaskQueueName,
	}

	// Execute workflow
	exec, err := c.temporal.ExecuteWorkflow(ctx, workflowOptions, api.ListExchangesWorkflowName, params)
	if err != nil {
		return api.ListExchangesResults{}, err
	}

	// Get result and return
	err = exec.Get(ctx, &res)
	return res, err
}