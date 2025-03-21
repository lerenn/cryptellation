package wfclient

import (
	api "github.com/lerenn/cryptellation/v1/api/worker/go"
	"github.com/lerenn/cryptellation/v1/clients/go/worker/wfclient/raw"
	"go.temporal.io/sdk/workflow"
)

func (c client) ListCandlesticks(
	ctx workflow.Context,
	params api.ListCandlesticksWorkflowParams,
	childWorkflowOptions *workflow.ChildWorkflowOptions,
) (result api.ListCandlesticksWorkflowResults, err error) {
	// TODO(#51): Implement caching
	return raw.ListCandlesticks(ctx, params, childWorkflowOptions)
}
