package workflows

import (
	"github.com/lerenn/cryptellation/v1/api"
	"go.temporal.io/sdk/workflow"
)

func (wf *workflows) UnregisterFromTicksListeningWorkflow(
	ctx workflow.Context,
	params api.UnregisterFromTicksListeningWorkflowParams,
) (api.UnregisterFromTicksListeningWorkflowResults, error) {
	// TODO
	return api.UnregisterFromTicksListeningWorkflowResults{}, nil
}
