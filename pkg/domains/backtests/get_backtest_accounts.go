package backtests

import (
	"fmt"

	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/pkg/domains/backtests/activities/db"
	"go.temporal.io/sdk/workflow"
)

func (wf *workflows) GetBacktestAccountsWorkflow(
	ctx workflow.Context,
	params api.GetBacktestAccountsWorkflowParams,
) (api.GetBacktestAccountsWorkflowResults, error) {
	// Read backtest
	var readRes db.ReadBacktestActivityResults
	err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, db.DefaultActivityOptions()),
		wf.db.ReadBacktestActivity, db.ReadBacktestActivityParams{
			ID: params.BacktestID,
		}).Get(ctx, &readRes)
	if err != nil {
		return api.GetBacktestAccountsWorkflowResults{}, fmt.Errorf("read backtest from db: %w", err)
	}

	return api.GetBacktestAccountsWorkflowResults{
		Accounts: readRes.Backtest.Accounts,
	}, nil
}
