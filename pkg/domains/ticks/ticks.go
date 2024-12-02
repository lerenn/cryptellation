package ticks

import (
	"github.com/lerenn/cryptellation/v1/api"
	temporalclient "go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

type Ticks interface {
	Register(w worker.Worker, temporal temporalclient.Client)

	RegisterForTicksListeningWorkflow(
		ctx workflow.Context,
		params api.RegisterForTicksListeningWorkflowParams,
	) (api.RegisterForTicksListeningWorkflowResults, error)

	UnregisterFromTicksListeningWorkflow(
		ctx workflow.Context,
		params api.UnregisterFromTicksListeningWorkflowParams,
	) (api.UnregisterFromTicksListeningWorkflowResults, error)
}
