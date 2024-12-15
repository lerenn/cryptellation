package backtests

import (
	"fmt"

	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/pkg/activities"
	"github.com/lerenn/cryptellation/v1/pkg/domains/backtests/activities/db"
	"go.temporal.io/sdk/workflow"
)

func (wf *workflows) GetBacktestWorkflow(
	ctx workflow.Context,
	params api.GetBacktestWorkflowParams,
) (api.GetBacktestWorkflowResults, error) {
	// Read backtest
	var readRes db.ReadBacktestActivityResults
	err := workflow.ExecuteActivity(workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: activities.DBInteractionDefaultTimeout,
	}), wf.db.ReadBacktestActivity, db.ReadBacktestActivityParams{
		ID: params.BacktestID,
	}).Get(ctx, &readRes)
	if err != nil {
		return api.GetBacktestWorkflowResults{}, fmt.Errorf("read backtest from db: %w", err)
	}

	return api.GetBacktestWorkflowResults{
		Backtest: readRes.Backtest,
	}, nil
}
