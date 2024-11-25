package workflows

import (
	"github.com/lerenn/cryptellation/v1/api"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

// Exchanges is the exchanges domain.
type Exchanges interface {
	Register(w worker.Worker)

	ListExchanges(ctx workflow.Context, params api.ListExchangesParams) (api.ListExchangesResults, error)
}
