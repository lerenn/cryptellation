package forwardtests

import (
	api "github.com/lerenn/cryptellation/v1/api/worker/go"
	"github.com/lerenn/cryptellation/v1/pkg/domains/forwardtests/activities/db"
	"go.temporal.io/sdk/workflow"
)

// ListForwardtestsWorkflow lists the forwardtests present in the system.
func (wf *workflows) ListForwardtestsWorkflow(
	ctx workflow.Context,
	_ api.ListForwardtestsWorkflowParams,
) (api.ListForwardtestsWorkflowResults, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Listing forwardtests")

	// List forwardtests
	var res db.ListForwardtestsActivityResult
	err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, db.DefaultActivityOptions()),
		wf.db.ListForwardtestsActivity, db.ListForwardtestsActivityParams{}).Get(ctx, &res)
	if err != nil {
		logger.Error("Error listing forwardtests",
			"error", err.Error())
		return api.ListForwardtestsWorkflowResults{}, err
	}

	logger.Info("Listed forwardtests",
		"count", len(res.Forwardtests))
	return api.ListForwardtestsWorkflowResults{
		Forwardtests: res.Forwardtests,
	}, nil
}
