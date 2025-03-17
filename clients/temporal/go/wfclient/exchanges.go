package wfclient

import (
	api "github.com/lerenn/cryptellation/v1/api/worker/go"
	"github.com/lerenn/cryptellation/v1/clients/temporal/go/wfclient/raw"
	"go.temporal.io/sdk/workflow"
)

func (c client) GetExchange(
	ctx workflow.Context,
	params api.GetExchangeWorkflowParams,
	childWorkflowOptions *workflow.ChildWorkflowOptions,
) (result api.GetExchangeWorkflowResults, err error) {
	// TODO(#52): Implement caching
	return raw.GetExchange(ctx, params, childWorkflowOptions)
}
