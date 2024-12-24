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
	var a *Activities
	var res SignalWithStartActivityResults
	return workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: time.Second * 10,
		}),
		a.SignalWithStartActivity,
		params).Get(ctx, &res)
}

type (
	// SignalWithStartActivityParams is the params for the SignalWithStartActivity activity.
	SignalWithStartActivityParams struct {
		SignalName     string
		SignalParams   any
		WorkflowID     string
		WorkflowName   string
		WorkflowParams any
		TaskQueue      string
	}

	// SignalWithStartActivityResults is the results from the SignalWithStartActivity activity.
	SignalWithStartActivityResults struct{}
)

// SignalWithStartActivity is an activity that will signal a workflow and start
// if before if it does not exist.
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
