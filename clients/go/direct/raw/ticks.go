package raw

import (
	"context"

	"github.com/lerenn/cryptellation/v1/api"
	temporalclient "go.temporal.io/sdk/client"
)

// ListenToTicks listens to ticks.
func (c raw) ListenToTicks(
	ctx context.Context,
	registerParams api.RegisterForTicksListeningWorkflowParams,
) (res api.RegisterForTicksListeningWorkflowResults, err error) {
	// Execute register workflow
	exec, err := c.temporal.ExecuteWorkflow(ctx,
		temporalclient.StartWorkflowOptions{
			TaskQueue: api.WorkerTaskQueueName,
		},
		api.RegisterForTicksListeningWorkflowName,
		registerParams)
	if err != nil {
		return api.RegisterForTicksListeningWorkflowResults{}, err
	}

	// Get result and return
	err = exec.Get(ctx, &res)
	return res, err
}
