package backtests

import (
	"fmt"

	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/pkg/domains/backtests/activities/db"
	"github.com/lerenn/cryptellation/v1/pkg/models/backtest"
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
	err = workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, db.DefaultActivityOptions()),
		wf.db.CreateBacktestActivity, db.CreateBacktestActivityParams{
			Backtest: bt,
		}).Get(ctx, &dbRes)
	if err != nil {
		return api.CreateBacktestWorkflowResults{}, fmt.Errorf("adding backtest to db: %w", err)
	}

	return api.CreateBacktestWorkflowResults{
		ID: bt.ID,
	}, nil
}
