package client

import (
	"context"

	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/pkg/config"
	temporalclient "go.temporal.io/sdk/client"
)

type client struct {
	temporal temporalclient.Client
}

// New creates a new client to execute temporal workflows.
// temporalConfig is the optional configuration to use for the temporal client.
func New(temporalConfig ...config.Temporal) (Client, error) {
	var t config.Temporal

	if len(temporalConfig) > 0 {
		t = temporalConfig[0]
	}

	// Load temporal configuration
	t = config.LoadTemporal(&t)
	if err := t.Validate(); err != nil {
		return client{}, err
	}

	// Load temporal client
	c, err := temporalclient.Dial(temporalclient.Options{
		HostPort: t.Address,
	})
	if err != nil {
		return client{}, err
	}

	return &client{temporal: c}, nil
}

// Info calls the service info.
func (c client) Info(ctx context.Context) (res api.ServiceInfoWorkflowResult, err error) {
	workflowOptions := temporalclient.StartWorkflowOptions{
		TaskQueue: api.WorkerTaskQueueName,
	}

	// Execute workflow
	exec, err := c.temporal.ExecuteWorkflow(ctx, workflowOptions, api.ServiceInfoWorkflowName)
	if err != nil {
		return api.ServiceInfoWorkflowResult{}, err
	}

	// Get result and return
	err = exec.Get(ctx, &res)
	return res, err
}

// Close closes the client.
func (c client) Close(ctx context.Context) {
	c.temporal.Close()
}
