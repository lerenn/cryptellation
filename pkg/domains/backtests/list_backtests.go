package backtests

import (
	"fmt"

	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/pkg/activities"
	"github.com/lerenn/cryptellation/v1/pkg/domains/backtests/activities/db"
	"go.temporal.io/sdk/workflow"
)

func (wf *workflows) ListBacktestsWorkflow(
	ctx workflow.Context,
	_ api.ListBacktestsWorkflowParams,
) (api.ListBacktestsWorkflowResults, error) {
	// Execute activity for listing backtests
	var dbRes db.ListBacktestsActivityResults
	err := workflow.ExecuteActivity(workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: activities.DBInteractionDefaultTimeout,
	}), wf.db.ListBacktestsActivity, db.ListBacktestsActivityParams{}).Get(ctx, &dbRes)
	if err != nil {
		return api.ListBacktestsWorkflowResults{}, fmt.Errorf("adding backtest to db: %w", err)
	}

	return api.ListBacktestsWorkflowResults{
		Backtests: dbRes.Backtests,
	}, nil
}
