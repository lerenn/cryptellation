package client

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/pkg/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

// Backtest is a local representation of a backtest running on the Cryptellation
// API.
type Backtest struct {
	ID uuid.UUID

	cryptellation client
}

// NewBacktest creates a new backtest.
func (c client) NewBacktest(
	ctx context.Context,
	params api.CreateBacktestWorkflowParams,
) (Backtest, error) {
	res, err := c.Raw.CreateBacktest(ctx, params)
	return Backtest{
		ID:            res.ID,
		cryptellation: c,
	}, err
}

// Run starts the backtest on Cryptellation API.
func (b Backtest) Run(ctx context.Context, robot Robot) error {
	// Create variables
	tq := fmt.Sprintf("CryptellationRunBacktest-%s", uuid.New().String())
	workflowName := tq

	// Create temporary worker
	w := worker.New(b.cryptellation.Temporal(), tq, worker.Options{})

	// Register OnInitCallback
	onInitCallbackWorkflowName := fmt.Sprintf("%s-OnInit", workflowName)
	w.RegisterWorkflowWithOptions(robot.OnInit, workflow.RegisterOptions{
		Name: onInitCallbackWorkflowName,
	})

	// Register OnNewPricesCallback
	onNewPricesCallbackWorkflowName := fmt.Sprintf("%s-OnPrice", workflowName)
	w.RegisterWorkflowWithOptions(robot.OnNewPrices, workflow.RegisterOptions{
		Name: onNewPricesCallbackWorkflowName,
	})

	// Register OnExitCallback
	onExitCallbackWorkflowName := fmt.Sprintf("%s-OnExit", workflowName)
	w.RegisterWorkflowWithOptions(robot.OnExit, workflow.RegisterOptions{
		Name: onExitCallbackWorkflowName,
	})

	// Start worker
	go func() {
		if err := w.Run(nil); err != nil {
			panic(err) // TODO: Handle error by returning it if there is an error
		}
	}()
	defer w.Stop()

	_, err := b.cryptellation.Raw.RunBacktest(ctx, api.RunBacktestWorkflowParams{
		BacktestID: b.ID,
		Callbacks: api.Callbacks{
			OnInitCallback: temporal.CallbackWorkflow{
				Name:          onInitCallbackWorkflowName,
				TaskQueueName: tq,
			},
			OnNewPricesCallback: temporal.CallbackWorkflow{
				Name:          onNewPricesCallbackWorkflowName,
				TaskQueueName: tq,
			},
			OnExitCallback: temporal.CallbackWorkflow{
				Name:          onExitCallbackWorkflowName,
				TaskQueueName: tq,
			},
		},
	})
	return err
}
