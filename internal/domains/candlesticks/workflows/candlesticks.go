package workflows

import (
	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/internal/domains/candlesticks/activities/db"
	"github.com/lerenn/cryptellation/v1/internal/domains/candlesticks/activities/exchanges"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

type candlesticks struct {
	db        db.Interface
	exchanges exchanges.Interface
}

// New creates a new candlesticks workflows.
func New(db db.Interface, exchanges exchanges.Interface) Candlesticks {
	if db == nil {
		panic("nil db")
	}

	if exchanges == nil {
		panic("nil exchanges")
	}

	return &candlesticks{
		db:        db,
		exchanges: exchanges,
	}
}

// Register registers the candlesticks workflows to the worker.
func (c *candlesticks) Register(w worker.Worker) {
	w.RegisterWorkflowWithOptions(c.ListCandlesticks, workflow.RegisterOptions{
		Name: api.ListCandlesticksWorkflowName,
	})
}
