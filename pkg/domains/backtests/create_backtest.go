package backtests

import (
	"fmt"
	"time"

	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/pkg/activities"
	"github.com/lerenn/cryptellation/v1/pkg/domains/backtests/activities/db"
	"github.com/lerenn/cryptellation/v1/pkg/models/backtest"
	temporalutils "github.com/lerenn/cryptellation/v1/pkg/temporal"
	"go.temporal.io/sdk/workflow"
)

// CreateBacktestWorkflow creates a new backtest and starts a workflow for running it.
func (wf *workflows) CreateBacktestWorkflow(
	ctx workflow.Context,
	params api.CreateBacktestWorkflowParams,
) (api.CreateBacktestWorkflowResults, error) {
	// Create backtest
	bt, err := backtest.New(params.BacktestParameters)
	if err != nil {
		return api.CreateBacktestWorkflowResults{}, fmt.Errorf("creating a new backtest from request: %w", err)
	}

	// Save it to DB
	var dbRes db.CreateBacktestActivityResults
	err = workflow.ExecuteActivity(workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: activities.DBInteractionDefaultTimeout,
	}), wf.db.CreateBacktestActivity, db.CreateBacktestActivityParams{
		Backtest: bt,
	}).Get(ctx, &dbRes)
	if err != nil {
		return api.CreateBacktestWorkflowResults{}, fmt.Errorf("adding backtest to db: %w", err)
	}

	return api.CreateBacktestWorkflowResults{
		ID: dbRes.ID,
	}, nil
}

func initBacktestCallback(
	ctx workflow.Context,
	onInitCallback temporalutils.CallbackWorkflow,
) error {
	// Options
	opts := workflow.ChildWorkflowOptions{
		TaskQueue:                onInitCallback.TaskQueueName, // Execute in the client queue
		WorkflowExecutionTimeout: time.Second * 30,             // Timeout if the child workflow does not complete
	}

	// Check if the timeout is set
	if onInitCallback.ExecutionTimeout > 0 {
		opts.WorkflowExecutionTimeout = onInitCallback.ExecutionTimeout
	}

	// Run a new child workflow
	ctx = workflow.WithChildOptions(ctx, opts)
	if err := workflow.ExecuteChildWorkflow(
		ctx, onInitCallback.Name, api.OnInitCallbackWorkflowParams{},
	).Get(ctx, nil); err != nil {
		return fmt.Errorf("starting new onInitCallback child workflow: %w", err)
	}

	return nil
}
