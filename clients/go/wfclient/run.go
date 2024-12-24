package wfclient

import (
	"errors"

	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/clients/go/wfclient/raw"
	"github.com/lerenn/cryptellation/v1/pkg/run"
	"go.temporal.io/sdk/workflow"
)

var (
	// ErrNotImplemented is returned when the function is not implemented.
	ErrNotImplemented = errors.New("not implemented")
)

// SubscribeToPriceParams is the parameters to subscribe to price updates.
type SubscribeToPriceParams struct {
	Run      run.Context
	Exchange string
	Pair     string
}

// SubscribeToPrice subscribes to specific price updates.
func (c client) SubscribeToPrice(ctx workflow.Context, params SubscribeToPriceParams) error {
	childWorkflowOptions := workflow.ChildWorkflowOptions{
		TaskQueue: params.Run.TaskQueue,
	}
	ctx = workflow.WithChildOptions(ctx, childWorkflowOptions)

	switch params.Run.Mode {
	case run.ModeBacktest:
		_, err := raw.SubscribeToBacktestPrice(ctx, api.SubscribeToBacktestPriceWorkflowParams{
			BacktestID: params.Run.ID,
			Exchange:   params.Exchange,
			Pair:       params.Pair,
		})
		return err
	case run.ModeForwardtest, run.ModeLive:
		return ErrNotImplemented
	default:
		return run.ErrInvalidMode
	}
}
