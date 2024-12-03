package workflows

import (
	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/pkg/domains/ticks"
	"github.com/lerenn/cryptellation/v1/pkg/domains/ticks/activities/exchanges"
	"github.com/lerenn/cryptellation/v1/pkg/domains/ticks/internal"
	temporalclient "go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

type workflows struct {
	exchanges exchanges.Exchanges
}

// New creates a new ticks workflows.
func New(exchanges exchanges.Exchanges) ticks.Ticks {
	// Check that the exchanges are valid
	if exchanges == nil {
		panic("nil exchanges")
	}

	return &workflows{
		exchanges: exchanges,
	}
}

func (wf *workflows) Register(w worker.Worker, temporalClient temporalclient.Client) {
	// Register internal workflows
	w.RegisterWorkflowWithOptions(wf.TicksSentryWorkflow, workflow.RegisterOptions{
		Name: internal.TicksSentryWorkflowName,
	})

	// Register public workflows
	w.RegisterWorkflowWithOptions(wf.RegisterForTicksListeningWorkflow, workflow.RegisterOptions{
		Name: api.RegisterForTicksListeningWorkflowName,
	})
	w.RegisterWorkflowWithOptions(wf.UnregisterFromTicksListeningWorkflow, workflow.RegisterOptions{
		Name: api.UnregisterFromTicksListeningWorkflowName,
	})
}
