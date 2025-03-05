package raw

import (
	"context"

	"github.com/lerenn/cryptellation/v1/api"
	temporalclient "go.temporal.io/sdk/client"
)

func (c raw) CreateForwardtest(
	ctx context.Context,
	params api.CreateForwardtestWorkflowParams,
) (api.CreateForwardtestWorkflowResults, error) {
	workflowOptions := temporalclient.StartWorkflowOptions{
		TaskQueue: api.WorkerTaskQueueName,
	}

	// Execute workflow
	exec, err := c.temporal.ExecuteWorkflow(ctx, workflowOptions, api.CreateForwardtestWorkflowName, params)
	if err != nil {
		return api.CreateForwardtestWorkflowResults{}, err
	}

	// Get result and return
	var res api.CreateForwardtestWorkflowResults
	err = exec.Get(ctx, &res)

	return res, err
}

func (c raw) GetForwardtestStatus(
	ctx context.Context,
	params api.GetForwardtestStatusWorkflowParams,
) (api.GetForwardtestStatusWorkflowResults, error) {
	workflowOptions := temporalclient.StartWorkflowOptions{
		TaskQueue: api.WorkerTaskQueueName,
	}

	// Execute workflow
	exec, err := c.temporal.ExecuteWorkflow(ctx, workflowOptions, api.GetForwardtestStatusWorkflowName, params)
	if err != nil {
		return api.GetForwardtestStatusWorkflowResults{}, err
	}

	// Get result and return
	var res api.GetForwardtestStatusWorkflowResults
	err = exec.Get(ctx, &res)

	return res, err
}

func (c raw) ListForwardtests(
	ctx context.Context,
	params api.ListForwardtestsWorkflowParams,
) (api.ListForwardtestsWorkflowResults, error) {
	workflowOptions := temporalclient.StartWorkflowOptions{
		TaskQueue: api.WorkerTaskQueueName,
	}

	// Execute workflow
	exec, err := c.temporal.ExecuteWorkflow(ctx, workflowOptions, api.ListForwardtestsWorkflowName, params)
	if err != nil {
		return api.ListForwardtestsWorkflowResults{}, err
	}

	// Get result and return
	var res api.ListForwardtestsWorkflowResults
	err = exec.Get(ctx, &res)

	return res, err
}

func (c raw) CreateForwardtestOrder(
	ctx context.Context,
	params api.CreateForwardtestOrderWorkflowParams,
) (api.CreateForwardtestOrderWorkflowResults, error) {
	workflowOptions := temporalclient.StartWorkflowOptions{
		TaskQueue: api.WorkerTaskQueueName,
	}

	// Execute workflow
	exec, err := c.temporal.ExecuteWorkflow(ctx, workflowOptions, api.CreateForwardtestOrderWorkflowName, params)
	if err != nil {
		return api.CreateForwardtestOrderWorkflowResults{}, err
	}

	// Get result and return
	var res api.CreateForwardtestOrderWorkflowResults
	err = exec.Get(ctx, &res)

	return res, err
}

func (c raw) ListForwardtestAccounts(
	ctx context.Context,
	params api.ListForwardtestAccountsWorkflowParams,
) (api.ListForwardtestAccountsWorkflowResults, error) {
	workflowOptions := temporalclient.StartWorkflowOptions{
		TaskQueue: api.WorkerTaskQueueName,
	}

	// Execute workflow
	exec, err := c.temporal.ExecuteWorkflow(ctx, workflowOptions, api.ListForwardtestAccountsWorkflowName, params)
	if err != nil {
		return api.ListForwardtestAccountsWorkflowResults{}, err
	}

	// Get result and return
	var res api.ListForwardtestAccountsWorkflowResults
	err = exec.Get(ctx, &res)

	return res, err
}
