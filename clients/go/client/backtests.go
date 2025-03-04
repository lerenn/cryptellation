package client

import (
	"context"

	"github.com/lerenn/cryptellation/v1/api"
)

// NewBacktest creates a new backtest.
func (c client) NewBacktest(
	ctx context.Context,
	params api.CreateBacktestWorkflowParams,
) (Backtest, error) {
	res, err := c.Client.CreateBacktest(ctx, params)
	return Backtest{
		ID:            res.ID,
		cryptellation: c,
	}, err
}

func (c client) GetBacktest(
	ctx context.Context,
	params api.GetBacktestWorkflowParams,
) (Backtest, error) {
	res, err := c.Client.GetBacktest(ctx, params)
	if err != nil {
		return Backtest{}, err
	}

	return Backtest{
		ID:            res.Backtest.ID,
		cryptellation: c,
	}, nil
}

func (c client) ListBacktests(
	ctx context.Context,
	params api.ListBacktestsWorkflowParams,
) ([]Backtest, error) {
	res, err := c.Client.ListBacktests(ctx, params)
	if err != nil {
		return nil, err
	}

	backtests := make([]Backtest, len(res.Backtests))
	for i, bt := range res.Backtests {
		backtests[i] = Backtest{
			ID:            bt.ID,
			cryptellation: c,
		}
	}

	return backtests, nil
}
