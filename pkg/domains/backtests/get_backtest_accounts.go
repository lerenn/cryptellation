package backtests

import (
	"fmt"

	"github.com/lerenn/cryptellation/v1/api"
	"go.temporal.io/sdk/workflow"
)

func (wf *workflows) GetBacktestAccountsWorkflow(
	ctx workflow.Context,
	params api.GetBacktestAccountsWorkflowParams,
) (api.GetBacktestAccountsWorkflowResults, error) {
	// Read backtest
	bt, err := wf.readBacktestFromDB(ctx, params.BacktestID)
	if err != nil {
		return api.GetBacktestAccountsWorkflowResults{}, fmt.Errorf("read backtest from db: %w", err)
	}

	return api.GetBacktestAccountsWorkflowResults{
		Accounts: bt.Accounts,
	}, nil
}
