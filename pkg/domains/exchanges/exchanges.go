package exchanges

import (
	"github.com/lerenn/cryptellation/v1/api"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

// Exchanges is the exchanges domain.
type Exchanges interface {
	Register(w worker.Worker)

	GetExchangeWorkflow(
		ctx workflow.Context,
		params api.GetExchangeWorkflowParams,
	) (api.GetExchangeWorkflowResults, error)

	ListExchangesWorkflow(
		ctx workflow.Context,
		params api.ListExchangesWorkflowParams,
	) (api.ListExchangesWorkflowResults, error)
}
