package ticks

import (
	"github.com/lerenn/cryptellation/v1/api"
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
