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

const DefaultExpirationDuration = time.Hour

func (e exchanges) ListExchanges(
	ctx workflow.Context,
	params api.ListExchangesParams,
) (api.ListExchangesResults, error) {
	// Log the start of the workflow
	workflow.GetLogger(ctx).Info(
		"Requested exchanges started",
		"names", params.Names)

	// Set activities params
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	// Read exchanges from database
	var dbRes db.ReadExchangesResult
	err := workflow.ExecuteActivity(ctx, e.db.ReadExchanges, db.ReadExchangesParams{
		Names: params.Names,
	}).Get(ctx, &dbRes)
	if err != nil {
		return api.ListExchangesResults{}, fmt.Errorf("handling exchanges from db reading: %w", err)
	}

	// Get the exchanges to synchronize
	toSync, err := exchange.GetExpiredExchangesNames(params.Names, dbRes.Exchanges, DefaultExpirationDuration)
	if err != nil {
		return api.ListExchangesResults{}, fmt.Errorf("determining exchanges to synchronize: %w", err)
	}

	// Get the exchanges from the services
	synced, err := e.getExchangeFromServices(ctx, toSync...)
	if err != nil {
		return api.ListExchangesResults{}, err
	}

	// Upsert the exchanges
	err = e.upsertExchanges(ctx, dbRes.Exchanges, synced)
	if err != nil {
		return api.ListExchangesResults{}, err
	}

	// Return the exchanges
	mappedExchanges := exchange.ArrayToMap(dbRes.Exchanges)
	for _, exch := range synced {
		mappedExchanges[exch.Name] = exch
	}
	return api.ListExchangesResults{
		List: exchange.MapToArray(mappedExchanges),
	}, nil
}

func (e exchanges) getExchangeFromServices(ctx workflow.Context, toSync ...string) ([]exchange.Exchange, error) {
	synced := make([]exchange.Exchange, 0, len(toSync))
	for _, name := range toSync {
		var r exchangesactivity.GetExchangeInfoResult
		err := workflow.ExecuteActivity(ctx, e.exchanges.GetExchangeInfo, exchangesactivity.GetExchangeInfoParams{
			Name: name,
		}).Get(ctx, &r)
		if err != nil {
			return nil, err
		}

		synced = append(synced, r.Exchange)
	}

	return synced, nil
}

func (e exchanges) upsertExchanges(ctx workflow.Context, dbExchanges, toUpsert []exchange.Exchange) error {
	toCreate := make([]exchange.Exchange, 0, len(toUpsert))
	toUpdate := make([]exchange.Exchange, 0, len(toUpsert))
	mappedDbExchanges := exchange.ArrayToMap(dbExchanges)
	for _, exch := range toUpsert {
		if _, ok := mappedDbExchanges[exch.Name]; ok {
			toUpdate = append(toUpdate, exch)
		} else {
			toCreate = append(toCreate, exch)
		}
	}

	if len(toCreate) > 0 {
		if err := workflow.ExecuteActivity(ctx, e.db.CreateExchanges, db.CreateExchangesParams{
			Exchanges: toCreate,
		}).Get(ctx, nil); err != nil {
			return err
		}
	}

	if len(toUpdate) > 0 {
		if err := workflow.ExecuteActivity(ctx, e.db.UpdateExchanges, db.UpdateExchangesParams{
			Exchanges: toUpdate,
		}).Get(ctx, nil); err != nil {
			return err
		}
	}

	return nil
}
