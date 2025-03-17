package client

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	api "github.com/lerenn/cryptellation/v1/api/worker/go"
	"github.com/lerenn/cryptellation/v1/clients/worker/go/bot"
	"github.com/lerenn/cryptellation/v1/pkg/run"
	"go.temporal.io/sdk/worker"
)

// Backtest is a local representation of a backtest running on the Cryptellation API.
type Backtest struct {
	ID uuid.UUID

	cryptellation client
}

// Run starts the backtest on Cryptellation API.
func (bt *Backtest) Run(ctx context.Context, b bot.Bot) error {
	// TODO(#49): get worker from parameters instead of creating a new one

	// Create temporary worker
	tq := fmt.Sprintf("%s-%s", run.ModeBacktest.String(), bt.ID.String())
	w := worker.New(bt.cryptellation.Temporal(), tq, worker.Options{})

	// Register workflows
	cbs := bot.RegisterWorkflows(w, tq, bt.ID, b)

	// Start worker
	go func() {
		if err := w.Run(nil); err != nil {
			panic(err)
		}
	}()
	defer w.Stop()

	_, err := bt.cryptellation.Client.RunBacktest(ctx, api.RunBacktestWorkflowParams{
		BacktestID: bt.ID,
		Callbacks:  cbs,
	})
	return err
}
