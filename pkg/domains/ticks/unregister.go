package ticks

import (
	api "github.com/lerenn/cryptellation/v1/api/worker/go"
	"go.temporal.io/sdk/workflow"
)

// UnregisterFromTicksListeningWorkflow will unregister a workflow from listening to ticks.
func (wf *workflows) UnregisterFromTicksListeningWorkflow(
	_ workflow.Context,
	_ api.UnregisterFromTicksListeningWorkflowParams,
) (api.UnregisterFromTicksListeningWorkflowResults, error) {
	// TODO(#68): Implement UnregisterFromTicksListeningWorkflow
	return api.UnregisterFromTicksListeningWorkflowResults{}, nil
}
