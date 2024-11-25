package workflows

import (
	"time"

	"github.com/lerenn/cryptellation/v1/api"
	exchangesactivity "github.com/lerenn/cryptellation/v1/internal/domains/exchanges/activities/exchanges"
	"go.temporal.io/sdk/workflow"
)

// ListExchanges will list the exchanges.
func (e exchanges) ListExchanges(
	ctx workflow.Context,
	_ api.ListExchangesParams,
) (api.ListExchangesResults, error) {
	// Set activities params
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	// Get the list of exchanges from the services
	var r exchangesactivity.ListExchangesNamesResult
	err := workflow.ExecuteActivity(ctx, e.exchanges.ListExchangesNames,
		exchangesactivity.ListExchangesNamesParams{}).Get(ctx, &r)
	if err != nil {
		return api.ListExchangesResults{}, err
	}

	// Return the result
	return api.ListExchangesResults{
		List: r.List,
	}, nil
}
