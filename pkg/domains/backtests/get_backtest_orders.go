package backtests

import (
	"fmt"

	"github.com/lerenn/cryptellation/v1/api"
	"go.temporal.io/sdk/workflow"
)

func (wf *workflows) GetBacktestOrdersWorkflow(
	ctx workflow.Context,
	params api.GetBacktestOrdersWorkflowParams,
) (api.GetBacktestOrdersWorkflowResults, error) {
	// Read backtest
	bt, err := wf.readBacktestFromDB(ctx, params.BacktestID)
	if err != nil {
		return api.GetBacktestOrdersWorkflowResults{}, fmt.Errorf("read backtest from db: %w", err)
	}

	return api.GetBacktestOrdersWorkflowResults{
		Orders: bt.Orders,
	}, nil
}
