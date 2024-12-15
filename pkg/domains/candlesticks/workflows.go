package candlesticks

import (
	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/pkg/domains/candlesticks/activities/db"
	"github.com/lerenn/cryptellation/v1/pkg/domains/candlesticks/activities/exchanges"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

type workflows struct {
	db        db.DB
	exchanges exchanges.Exchanges
}

// New creates a new candlesticks workflows.
func New(db db.DB, exchanges exchanges.Exchanges) Candlesticks {
	if db == nil {
		panic("nil db")
	}

	if exchanges == nil {
		panic("nil exchanges")
	}

	return &workflows{
		db:        db,
		exchanges: exchanges,
	}
}

// Register registers the candlesticks workflows to the worker.
func (wf *workflows) Register(w worker.Worker) {
	w.RegisterWorkflowWithOptions(wf.ListCandlesticksWorkflow, workflow.RegisterOptions{
		Name: api.ListCandlesticksWorkflowName,
	})
}
