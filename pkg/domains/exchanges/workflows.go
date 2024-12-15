package exchanges

import (
	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/pkg/domains/exchanges/activities/db"
	exchangesadapter "github.com/lerenn/cryptellation/v1/pkg/domains/exchanges/activities/exchanges"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

type workflows struct {
	db        db.DB
	exchanges exchangesadapter.Exchanges
}

// New creates a new exchanges workflows.
func New(db db.DB, exchs exchangesadapter.Exchanges) Exchanges {
	if db == nil {
		panic("nil db")
	}

	if exchs == nil {
		panic("nil exchanges")
	}

	return &workflows{
		db:        db,
		exchanges: exchs,
	}
}

// Register registers the candlesticks workflows to the worker.
func (wf *workflows) Register(w worker.Worker) {
	w.RegisterWorkflowWithOptions(wf.GetExchangeWorkflow, workflow.RegisterOptions{
		Name: api.GetExchangeWorkflowName,
	})
	w.RegisterWorkflowWithOptions(wf.ListExchangesWorkflow, workflow.RegisterOptions{
		Name: api.ListExchangesWorkflowName,
	})
}
