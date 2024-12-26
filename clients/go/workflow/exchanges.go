package wfclient

import (
	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/clients/go/workflow/raw"
	"go.temporal.io/sdk/workflow"
)

func (c client) GetExchange(
	ctx workflow.Context,
	params api.GetExchangeWorkflowParams,
	childWorkflowOptions *workflow.ChildWorkflowOptions,
) (result api.GetExchangeWorkflowResults, err error) {
	// TODO: Implement caching
	return raw.GetExchange(ctx, params, childWorkflowOptions)
}
