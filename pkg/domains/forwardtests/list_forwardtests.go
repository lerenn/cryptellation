package forwardtests

import (
	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/pkg/domains/forwardtests/activities/db"
	"go.temporal.io/sdk/workflow"
)

// ListForwardTestsWorkflow lists the forwardtests present in the system.
func (wf *workflows) ListForwardTestsWorkflow(
	ctx workflow.Context,
	_ api.ListForwardTestsWorkflowParams,
) (api.ListForwardTestsWorkflowResults, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Listing forwardtests")

	// List forwardtests
	var res db.ListForwardTestsActivityResult
	err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, db.DefaultActivityOptions()),
		wf.db.ListForwardTestsActivity, db.ListForwardTestsActivityParams{}).Get(ctx, &res)
	if err != nil {
		logger.Error("Error listing forwardtests",
			"error", err.Error())
		return api.ListForwardTestsWorkflowResults{}, err
	}

	logger.Info("Listed forwardtests",
		"count", len(res.ForwardTests))
	return api.ListForwardTestsWorkflowResults{
		ForwardTests: res.ForwardTests,
	}, nil
}
