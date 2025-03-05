package raw

import (
	"context"

	"github.com/lerenn/cryptellation/v1/api"
	temporalclient "go.temporal.io/sdk/client"
)

func (c raw) CreateBacktest(
	ctx context.Context,
	params api.CreateBacktestWorkflowParams,
) (api.CreateBacktestWorkflowResults, error) {
	workflowOptions := temporalclient.StartWorkflowOptions{
		TaskQueue: api.WorkerTaskQueueName,
	}

	// Execute workflow
	exec, err := c.temporal.ExecuteWorkflow(ctx, workflowOptions, api.CreateBacktestWorkflowName, params)
	if err != nil {
		return api.CreateBacktestWorkflowResults{}, err
	}

	// Get result and return
	var res api.CreateBacktestWorkflowResults
	err = exec.Get(ctx, &res)

	return res, err
}

func (c raw) RunBacktest(
	ctx context.Context,
	params api.RunBacktestWorkflowParams,
) (api.RunBacktestWorkflowResults, error) {
	workflowOptions := temporalclient.StartWorkflowOptions{
		TaskQueue: api.WorkerTaskQueueName,
	}

	// Execute workflow
	exec, err := c.temporal.ExecuteWorkflow(ctx, workflowOptions, api.RunBacktestWorkflowName, params)
	if err != nil {
		return api.RunBacktestWorkflowResults{}, err
	}

	// Get result and return
	var res api.RunBacktestWorkflowResults
	err = exec.Get(ctx, &res)

	return res, err
}

func (c raw) SubscribeToBacktestPrice(
	ctx context.Context,
	params api.SubscribeToBacktestPriceWorkflowParams,
) (api.SubscribeToBacktestPriceWorkflowResults, error) {
	workflowOptions := temporalclient.StartWorkflowOptions{
		TaskQueue: api.WorkerTaskQueueName,
	}

	// Execute workflow
	exec, err := c.temporal.ExecuteWorkflow(ctx, workflowOptions, api.SubscribeToBacktestPriceWorkflowName, params)
	if err != nil {
		return api.SubscribeToBacktestPriceWorkflowResults{}, err
	}

	// Get result and return
	var res api.SubscribeToBacktestPriceWorkflowResults
	err = exec.Get(ctx, &res)

	return res, err
}

func (c raw) ListBacktests(
	ctx context.Context,
	params api.ListBacktestsWorkflowParams,
) (api.ListBacktestsWorkflowResults, error) {
	workflowOptions := temporalclient.StartWorkflowOptions{
		TaskQueue: api.WorkerTaskQueueName,
	}

	// Execute workflow
	exec, err := c.temporal.ExecuteWorkflow(ctx, workflowOptions, api.ListBacktestsWorkflowName, params)
	if err != nil {
		return api.ListBacktestsWorkflowResults{}, err
	}

	// Get result and return
	var res api.ListBacktestsWorkflowResults
	err = exec.Get(ctx, &res)

	return res, err
}

func (c raw) GetBacktest(
	ctx context.Context,
	params api.GetBacktestWorkflowParams,
) (api.GetBacktestWorkflowResults, error) {
	workflowOptions := temporalclient.StartWorkflowOptions{
		TaskQueue: api.WorkerTaskQueueName,
	}

	// Execute workflow
	exec, err := c.temporal.ExecuteWorkflow(ctx, workflowOptions, api.GetBacktestWorkflowName, params)
	if err != nil {
		return api.GetBacktestWorkflowResults{}, err
	}

	// Get result and return
	var res api.GetBacktestWorkflowResults
	err = exec.Get(ctx, &res)

	return res, err
}
