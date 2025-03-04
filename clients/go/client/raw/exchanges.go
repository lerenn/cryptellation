package raw

import (
	"context"

	"github.com/lerenn/cryptellation/v1/api"
	temporalclient "go.temporal.io/sdk/client"
)

// GetExchange calls the exchange get workflow.
func (c raw) GetExchange(
	ctx context.Context,
	params api.GetExchangeWorkflowParams,
) (res api.GetExchangeWorkflowResults, err error) {
	workflowOptions := temporalclient.StartWorkflowOptions{
		TaskQueue: api.WorkerTaskQueueName,
	}

	// Execute workflow
	exec, err := c.temporal.ExecuteWorkflow(ctx, workflowOptions, api.GetExchangeWorkflowName, params)
	if err != nil {
		return api.GetExchangeWorkflowResults{}, err
	}

	// Get result and return
	err = exec.Get(ctx, &res)
	return res, err
}

// ListExchanges calls the exchanges list workflow.
func (c raw) ListExchanges(
	ctx context.Context,
	params api.ListExchangesWorkflowParams,
) (res api.ListExchangesWorkflowResults, err error) {
	workflowOptions := temporalclient.StartWorkflowOptions{
		TaskQueue: api.WorkerTaskQueueName,
	}

	// Execute workflow
	exec, err := c.temporal.ExecuteWorkflow(ctx, workflowOptions, api.ListExchangesWorkflowName, params)
	if err != nil {
		return api.ListExchangesWorkflowResults{}, err
	}

	// Get result and return
	err = exec.Get(ctx, &res)
	return res, err
}
