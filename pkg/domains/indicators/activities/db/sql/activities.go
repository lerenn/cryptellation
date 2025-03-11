package sql

import (
	"context"
	"fmt"

	"github.com/lerenn/cryptellation/v1/pkg/activities"
	"github.com/lerenn/cryptellation/v1/pkg/config"
	"github.com/lerenn/cryptellation/v1/pkg/domains/indicators/activities/db"
	"github.com/lerenn/cryptellation/v1/pkg/domains/indicators/activities/db/sql/entities"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
)

var _ db.DB = (*Activities)(nil)

// Activities is a struct that contains all the methods to interact with the
// activities table in the database.
type Activities struct {
	client activities.PostGres
}

// New creates a new activities.
func New(ctx context.Context, c config.PostGres) (*Activities, error) {
	// Create embedded database access
	db, err := activities.NewPostGres(ctx, c)
	if err != nil {
		return nil, err
	}

	// Create a structure
	a := &Activities{
		client: db,
	}

	return a, nil
}

// Register registers the activities.
func (a *Activities) Register(w worker.Worker) {
	w.RegisterActivityWithOptions(
		a.ReadSMAActivity,
		activity.RegisterOptions{Name: db.ReadSMAActivityName},
	)
	w.RegisterActivityWithOptions(
		a.UpsertSMAActivity,
		activity.RegisterOptions{Name: db.UpsertSMAActivityName},
	)
}

// Reset will reset the database.
func (a *Activities) Reset(ctx context.Context) error {
	_, err := a.client.DB.ExecContext(ctx, "DELETE FROM indicators_sma")
	if err != nil {
		return fmt.Errorf("deleting indicators_sma rows: %w", err)
	}

	return nil
}

// ReadSMAActivity reads the SMA points from the database.
func (a *Activities) ReadSMAActivity(
	ctx context.Context,
	params db.ReadSMAActivityParams,
) (db.ReadSMAActivityResults, error) {
	// Query the SMA points
	rows, err := a.client.DB.QueryxContext(
		ctx,
		`SELECT *
		FROM indicators_sma
		WHERE exchange = $1 AND pair = $2 AND period = $3 AND period_number = $4 AND price_type = $5 AND time >= $6 AND time <= $7
		ORDER BY time ASC`,
		params.Exchange,
		params.Pair,
		params.Period,
		params.PeriodNumber,
		params.PriceType,
		params.Start.UTC(),
		params.End.UTC(),
	)
	if err != nil {
		return db.ReadSMAActivityResults{}, fmt.Errorf("querying SMA points: %w", err)
	}
	defer rows.Close()

	// Loop through the rows
	results := make([]entities.SimpleMovingAverage, 0)
	for rows.Next() {
		// Create the SMA point
		var point entities.SimpleMovingAverage
		err = rows.StructScan(&point)
		if err != nil {
			return db.ReadSMAActivityResults{}, fmt.Errorf("scanning SMA point: %w", err)
		}

		// Append the point
		results = append(results, point)
	}

	// To model list
	data, err := entities.FromEntityListToModelList(results)
	if err != nil {
		return db.ReadSMAActivityResults{}, fmt.Errorf("from entity list to model list: %w", err)
	}

	// Return the results
	return db.ReadSMAActivityResults{
		Data: data,
	}, nil
}

// UpsertSMAActivity upserts the SMA points in the database.
func (a *Activities) UpsertSMAActivity(
	ctx context.Context,
	params db.UpsertSMAActivityParams,
) (db.UpsertSMAActivityResults, error) {
	// Create entities
	ents := entities.FromModelListToEntityList(
		params.Exchange,
		params.Pair,
		params.Period,
		params.PeriodNumber,
		params.PriceType,
		params.TimeSerie)

	// Bulk insert the SMA
	_, err := a.client.DB.NamedExecContext(
		ctx,
		`INSERT INTO indicators_sma (exchange, pair, period, period_number, price_type, time, data)
		VALUES (:exchange, :pair, :period, :period_number, :price_type, :time, :data)
		ON CONFLICT (exchange, pair, period, period_number, price_type, time) DO UPDATE
		SET data = EXCLUDED.data`,
		entities.FromEntitiesToMap(ents),
	)
	if err != nil {
		return db.UpsertSMAActivityResults{}, fmt.Errorf("bulk inserting sma: %w", err)
	}

	return db.UpsertSMAActivityResults{}, nil
}
