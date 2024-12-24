package backtests

import (
	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/v1/pkg/domains/backtests/activities/db"
	"github.com/lerenn/cryptellation/v1/pkg/models/backtest"
	"go.temporal.io/sdk/workflow"
)

func (wf *workflows) readBacktestFromDB(ctx workflow.Context, id uuid.UUID) (backtest.Backtest, error) {
	var readRes db.ReadBacktestActivityResults
	err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, db.DefaultActivityOptions()),
		wf.db.ReadBacktestActivity, db.ReadBacktestActivityParams{
			ID: id,
		}).Get(ctx, &readRes)
	if err != nil {
		return backtest.Backtest{}, err
	}

	return readRes.Backtest, nil
}
