package backtests

import (
	"fmt"

	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/pkg/domains/backtests/activities/db"
	"go.temporal.io/sdk/workflow"
)

func (wf *workflows) GetBacktestOrdersWorkflow(
	ctx workflow.Context,
	params api.GetBacktestOrdersWorkflowParams,
) (api.GetBacktestOrdersWorkflowResults, error) {
	// Read backtest
	var readRes db.ReadBacktestActivityResults
	err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, db.DefaultActivityOptions()),
		wf.db.ReadBacktestActivity, db.ReadBacktestActivityParams{
			ID: params.BacktestID,
		}).Get(ctx, &readRes)
	if err != nil {
		return api.GetBacktestOrdersWorkflowResults{}, fmt.Errorf("read backtest from db: %w", err)
	}

	return api.GetBacktestOrdersWorkflowResults{
		Orders: readRes.Backtest.Orders,
	}, nil
}
