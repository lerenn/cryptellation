package workflows

import (
	"fmt"
	"time"

	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/internal/domains/exchanges/activities/db"
	exchangesactivity "github.com/lerenn/cryptellation/v1/internal/domains/exchanges/activities/exchanges"
	"github.com/lerenn/cryptellation/v1/pkg/models/exchange"
	"go.temporal.io/sdk/workflow"
)

// DefaultExpirationDuration is the default duration for an exchange to be considered outdated.
const DefaultExpirationDuration = time.Hour

// GetExchange will get a specific exchange.
func (e exchanges) GetExchange(ctx workflow.Context, params api.GetExchangeParams) (api.GetExchangeResults, error) {
	// Log the start of the workflow
	workflow.GetLogger(ctx).Info(
		"Requested exchange started",
		"name", params.Name)

	// Set activities params
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	// Read exchanges from database
	var dbRes db.ReadExchangesResult
	err := workflow.ExecuteActivity(ctx, e.db.ReadExchanges, db.ReadExchangesParams{
		Names: []string{params.Name},
	}).Get(ctx, &dbRes)
	if err != nil {
		return api.GetExchangeResults{}, fmt.Errorf("handling exchanges from db reading: %w", err)
	}

	// Check if the exchange is present and not expired
	if len(dbRes.Exchanges) > 0 && !dbRes.Exchanges[0].IsOutdated(DefaultExpirationDuration) {
		return api.GetExchangeResults{
			Exchange: dbRes.Exchanges[0],
		}, nil
	}

	// Get the exchange from the services
	var r exchangesactivity.GetExchangeInfoResult
	err = workflow.ExecuteActivity(ctx, e.exchanges.GetExchangeInfo, exchangesactivity.GetExchangeInfoParams{
		Name: params.Name,
	}).Get(ctx, &r)
	if err != nil {
		return api.GetExchangeResults{}, err
	}

	// Save the exchange in the database
	if len(dbRes.Exchanges) == 0 {
		if err := workflow.ExecuteActivity(ctx, e.db.CreateExchanges, db.CreateExchangesParams{
			Exchanges: []exchange.Exchange{r.Exchange},
		}).Get(ctx, nil); err != nil {
			return api.GetExchangeResults{}, err
		}
	} else {
		if err := workflow.ExecuteActivity(ctx, e.db.UpdateExchanges, db.UpdateExchangesParams{
			Exchanges: []exchange.Exchange{r.Exchange},
		}).Get(ctx, nil); err != nil {
			return api.GetExchangeResults{}, err
		}
	}

	// Return the result
	return api.GetExchangeResults{
		Exchange: r.Exchange,
	}, nil
}
