package workflows

import (
	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/internal/domains/exchanges/activities/db"
	exchangesadapter "github.com/lerenn/cryptellation/v1/internal/domains/exchanges/activities/exchanges"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

type exchanges struct {
	db        db.Interface
	exchanges exchangesadapter.Interface
}

// New creates a new exchanges workflows.
func New(db db.Interface, exchs exchangesadapter.Interface) Exchanges {
	if db == nil {
		panic("nil db")
	}

	if exchs == nil {
		panic("nil exchanges")
	}

	return &exchanges{
		db:        db,
		exchanges: exchs,
	}
}

// Register registers the candlesticks workflows to the worker.
func (e exchanges) Register(w worker.Worker) {
	w.RegisterWorkflowWithOptions(e.GetExchange, workflow.RegisterOptions{
		Name: api.GetExchangeWorkflowName,
	})
	w.RegisterWorkflowWithOptions(e.ListExchanges, workflow.RegisterOptions{
		Name: api.ListExchangesWorkflowName,
	})
}
