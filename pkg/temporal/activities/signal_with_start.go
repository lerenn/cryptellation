package activities

import (
	"context"
	"time"

	temporalclient "go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
)

// ExecuteSignalWithStart is a wrapper for the SignalWithStartActivity execution.
func ExecuteSignalWithStart(
	ctx workflow.Context,
	params SignalWithStartActivityParams,
) error {
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	var (
		a   *Activities
		res SignalWithStartActivityResults
	)
	return workflow.ExecuteActivity(
		ctx,
		a.SignalWithStartActivity,
		params).Get(ctx, &res)
}

type (
	SignalWithStartActivityParams struct {
		SignalName     string
		SignalParams   any
		WorkflowID     string
		WorkflowName   string
		WorkflowParams any
		TaskQueue      string
	}

	SignalWithStartActivityResults struct{}
)

func (a *Activities) SignalWithStartActivity(
	ctx context.Context,
	params SignalWithStartActivityParams,
) (SignalWithStartActivityResults, error) {
	_, err := a.temporal.SignalWithStartWorkflow(
		ctx,
		params.WorkflowID,
		params.SignalName,
		params.SignalParams,
		temporalclient.StartWorkflowOptions{TaskQueue: params.TaskQueue},
		params.WorkflowName,
		params.WorkflowParams,
	)
	return SignalWithStartActivityResults{}, err
}
