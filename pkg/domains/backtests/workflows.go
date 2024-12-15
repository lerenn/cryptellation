package backtests

import (
	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/pkg/domains/backtests/activities/db"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

type workflows struct {
	db db.DB
}

// New creates a new backtests workflows.
func New(db db.DB) Backtests {
	if db == nil {
		panic("nil db")
	}

	return &workflows{
		db: db,
	}
}

// Register registers the candlesticks workflows to the worker.
func (wf *workflows) Register(w worker.Worker) {
	w.RegisterWorkflowWithOptions(wf.CreateBacktestOrderWorkflow, workflow.RegisterOptions{
		Name: api.CreateBacktestOrderWorkflowName,
	})
	w.RegisterWorkflowWithOptions(wf.CreateBacktestWorkflow, workflow.RegisterOptions{
		Name: api.CreateBacktestWorkflowName,
	})
	w.RegisterWorkflowWithOptions(wf.GetBacktestAccountsWorkflow, workflow.RegisterOptions{
		Name: api.GetBacktestAccountsWorkflowName,
	})
	w.RegisterWorkflowWithOptions(wf.GetBacktestOrdersWorkflow, workflow.RegisterOptions{
		Name: api.GetBacktestOrdersWorkflowName,
	})
	w.RegisterWorkflowWithOptions(wf.GetBacktestWorkflow, workflow.RegisterOptions{
		Name: api.GetBacktestWorkflowName,
	})
	w.RegisterWorkflowWithOptions(wf.ListBacktestsWorkflow, workflow.RegisterOptions{
		Name: api.ListBacktestsWorkflowName,
	})
	w.RegisterWorkflowWithOptions(wf.RunBacktestWorkflow, workflow.RegisterOptions{
		Name: api.RunBacktestWorkflowName,
	})
}
