package workflow

import (
	"errors"

	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/clients/go/workflow/raw"
	"github.com/lerenn/cryptellation/v1/pkg/run"
	"go.temporal.io/sdk/workflow"
)

var (
	ErrNotImplemented = errors.New("not implemented")
)

type SubscribeToPriceParams struct {
	Run      run.Context
	Exchange string
	Pair     string
}

func SubscribeToPrice(ctx workflow.Context, params SubscribeToPriceParams) error {
	childWorkflowOptions := workflow.ChildWorkflowOptions{
		TaskQueue: params.Run.TaskQueue,
	}
	ctx = workflow.WithChildOptions(ctx, childWorkflowOptions)

	switch params.Run.Mode {
	case run.ModeBacktest:
		raw.SubscribeToBacktestPrice(ctx, api.SubscribeToBacktestPriceWorkflowParams{
			BacktestID: params.Run.ID,
			Exchange:   params.Exchange,
			Pair:       params.Pair,
		})
		return nil
	case run.ModeForwardtest, run.ModeLive:
		return ErrNotImplemented
	default:
		return run.ErrInvalidMode
	}
}
