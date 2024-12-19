package direct

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/clients/go/bot"
	"github.com/lerenn/cryptellation/v1/pkg/run"
	"go.temporal.io/sdk/worker"
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
	res, err := c.Client.CreateBacktest(ctx, params)
	return Backtest{
		ID:            res.ID,
		cryptellation: c,
	}, err
}

// Run starts the backtest on Cryptellation API.
func (bt *Backtest) Run(ctx context.Context, b bot.Bot) error {
	// Create temporary worker
	tq := fmt.Sprintf("%s-%s", run.ModeBacktest.String(), bt.ID.String())
	w := worker.New(bt.cryptellation.Temporal(), tq, worker.Options{})

	// Register workflows
	cbs := bot.RegisterWorkflows(w, tq, bt.ID, b)

	// Start worker
	go func() {
		if err := w.Run(nil); err != nil {
			panic(err) // TODO: Handle error by returning it if there is an error
		}
	}()
	defer w.Stop()

	_, err := bt.cryptellation.Client.RunBacktest(ctx, api.RunBacktestWorkflowParams{
		BacktestID: bt.ID,
		Callbacks:  cbs,
	})
	return err
}
