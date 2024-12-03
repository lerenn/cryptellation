package workflows

import (
	"time"

	"github.com/lerenn/cryptellation/v1/api"
	exchangesactivity "github.com/lerenn/cryptellation/v1/pkg/domains/exchanges/activities/exchanges"
	"go.temporal.io/sdk/workflow"
)

// ListExchangesWorkflow will list the exchanges.
func (wf *workflows) ListExchangesWorkflow(
	ctx workflow.Context,
	_ api.ListExchangesWorkflowParams,
) (api.ListExchangesWorkflowResults, error) {
	// Set activities params
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	// Get the list of exchanges from the services
	var r exchangesactivity.ListExchangesActivityResults
	err := workflow.ExecuteActivity(ctx, wf.exchanges.ListExchangesActivity,
		exchangesactivity.ListExchangesActivityParams{}).Get(ctx, &r)
	if err != nil {
		return api.ListExchangesWorkflowResults{}, err
	}

	// Return the result
	return api.ListExchangesWorkflowResults{
		List: r.List,
	}, nil
}
