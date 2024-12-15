package ticks

import (
	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/pkg/domains/ticks/activities/exchanges"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

type workflows struct {
	exchanges exchanges.Exchanges
}

// New creates a new ticks workflows.
func New(exchanges exchanges.Exchanges) Ticks {
	// Check that the exchanges are valid
	if exchanges == nil {
		panic("nil exchanges")
	}

	return &workflows{
		exchanges: exchanges,
	}
}

func (wf *workflows) Register(w worker.Worker) {
	// Private workflows
	w.RegisterWorkflowWithOptions(wf.ticksSentryWorkflow, workflow.RegisterOptions{
		Name: ticksSentryWorkflowName,
	})

	// Public workflows
	w.RegisterWorkflowWithOptions(wf.RegisterForTicksListeningWorkflow, workflow.RegisterOptions{
		Name: api.RegisterForTicksListeningWorkflowName,
	})
	w.RegisterWorkflowWithOptions(wf.UnregisterFromTicksListeningWorkflow, workflow.RegisterOptions{
		Name: api.UnregisterFromTicksListeningWorkflowName,
	})
}
