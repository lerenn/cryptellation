package exchanges

import (
	"fmt"
	"time"

	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/pkg/domains/exchanges/activities/db"
	"github.com/lerenn/cryptellation/v1/pkg/domains/exchanges/activities/exchanges"
	"github.com/lerenn/cryptellation/v1/pkg/models/exchange"
	"go.temporal.io/sdk/workflow"
)

// DefaultExpirationDuration is the default duration for an exchange to be considered outdated.
const DefaultExpirationDuration = time.Hour

// GetExchangeWorkflow will get a specific exchange.
func (wf *workflows) GetExchangeWorkflow(
	ctx workflow.Context,
	params api.GetExchangeWorkflowParams,
) (api.GetExchangeWorkflowResults, error) {
	// Log the start of the workflow
	workflow.GetLogger(ctx).Info(
		"Requested exchange started",
		"name", params.Name)

	// Read exchanges from database
	var dbRes db.ReadExchangesActivityResults
	err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, db.DefaultActivityOptions()),
		wf.db.ReadExchangesActivity, db.ReadExchangesActivityParams{
			Names: []string{params.Name},
		}).Get(ctx, &dbRes)
	if err != nil {
		return api.GetExchangeWorkflowResults{}, fmt.Errorf("handling exchanges from db reading: %w", err)
	}

	// Check if the exchange is present and not expired
	if len(dbRes.Exchanges) > 0 && !dbRes.Exchanges[0].IsOutdated(DefaultExpirationDuration) {
		return api.GetExchangeWorkflowResults{
			Exchange: dbRes.Exchanges[0],
		}, nil
	}

	// Get the exchange from the services
	var r exchanges.GetExchangeActivityResults
	err = workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, exchanges.DefaultActivityOptions()),
		wf.exchanges.GetExchangeActivity, exchanges.GetExchangeActivityParams{
			Name: params.Name,
		}).Get(ctx, &r)
	if err != nil {
		return api.GetExchangeWorkflowResults{}, err
	}

	// Save the exchange in the database
	if len(dbRes.Exchanges) == 0 {
		if err := workflow.ExecuteActivity(
			workflow.WithActivityOptions(ctx, db.DefaultActivityOptions()),
			wf.db.CreateExchangesActivity, db.CreateExchangesActivityParams{
				Exchanges: []exchange.Exchange{r.Exchange},
			}).Get(ctx, nil); err != nil {
			return api.GetExchangeWorkflowResults{}, err
		}
	} else {
		if err := workflow.ExecuteActivity(
			workflow.WithActivityOptions(ctx, db.DefaultActivityOptions()),
			wf.db.UpdateExchangesActivity, db.UpdateExchangesActivityParams{
				Exchanges: []exchange.Exchange{r.Exchange},
			}).Get(ctx, nil); err != nil {
			return api.GetExchangeWorkflowResults{}, err
		}
	}

	// Return the result
	return api.GetExchangeWorkflowResults{
		Exchange: r.Exchange,
	}, nil
}
