package wfclient

import (
	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/clients/go/workflow/raw"
	"go.temporal.io/sdk/workflow"
)

func (c client) ListCandlesticks(
	ctx workflow.Context,
	params api.ListCandlesticksWorkflowParams,
	childWorkflowOptions *workflow.ChildWorkflowOptions,
) (result api.ListCandlesticksWorkflowResults, err error) {
	// TODO: Implement caching
	return raw.ListCandlesticks(ctx, params, childWorkflowOptions)
}
