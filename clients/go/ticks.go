package client

import (
	"context"

	"github.com/lerenn/cryptellation/v1/api"
	temporalclient "go.temporal.io/sdk/client"
)

func (c client) ListenToTicks(
	ctx context.Context,
	params api.RegisterForTicksListeningWorkflowParams,
) (res api.RegisterForTicksListeningWorkflowResults, err error) {
	workflowOptions := temporalclient.StartWorkflowOptions{
		TaskQueue: api.WorkerTaskQueueName,
	}

	// Execute workflow
	exec, err := c.temporal.ExecuteWorkflow(ctx, workflowOptions, api.RegisterForTicksListeningWorkflowName, params)
	if err != nil {
		return api.RegisterForTicksListeningWorkflowResults{}, err
	}

	// Get result and return
	err = exec.Get(ctx, &res)
	return res, err
}
