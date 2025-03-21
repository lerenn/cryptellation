package ticks

import (
	api "github.com/lerenn/cryptellation/v1/api/worker/go"
	"github.com/lerenn/cryptellation/v1/clients/go/worker/wfclient"
	"github.com/lerenn/cryptellation/v1/pkg/domains/ticks/activities/exchanges"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

// Ticks is the ticks domain.
type Ticks interface {
	Register(w worker.Worker)

	RegisterForTicksListeningWorkflow(
		ctx workflow.Context,
		params api.RegisterForTicksListeningWorkflowParams,
	) (api.RegisterForTicksListeningWorkflowResults, error)

	UnregisterFromTicksListeningWorkflow(
		ctx workflow.Context,
		params api.UnregisterFromTicksListeningWorkflowParams,
	) (api.UnregisterFromTicksListeningWorkflowResults, error)
}

// Check that the workflows implements the Ticks interface.
var _ Ticks = &workflows{}

type workflows struct {
	exchanges     exchanges.Exchanges
	cryptellation wfclient.Client
}

// New creates a new ticks workflows.
func New(exchanges exchanges.Exchanges) Ticks {
	// Check that the exchanges are valid
	if exchanges == nil {
		panic("nil exchanges")
	}

	return &workflows{
		cryptellation: wfclient.NewClient(),
		exchanges:     exchanges,
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
