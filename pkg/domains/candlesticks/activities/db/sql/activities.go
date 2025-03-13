package sql

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lerenn/cryptellation/v1/pkg/activities"
	"github.com/lerenn/cryptellation/v1/pkg/config"
	"github.com/lerenn/cryptellation/v1/pkg/domains/candlesticks/activities/db"
	"github.com/lerenn/cryptellation/v1/pkg/domains/candlesticks/activities/db/sql/entities"
	"github.com/lerenn/cryptellation/v1/pkg/models/candlestick"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
)

var _ db.DB = (*Activities)(nil)

// Activities is a struct that contains all the methods to interact with the
// activities table in the database.
type Activities struct {
	client activities.SQL
}

// New creates a new activities.
func New(ctx context.Context, c config.SQL) (*Activities, error) {
	// Create embedded database access
	db, err := activities.NewSQL(ctx, c)
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
		a.CreateCandlesticksActivity,
		activity.RegisterOptions{Name: db.CreateCandlesticksActivityName},
	)
	w.RegisterActivityWithOptions(
		a.ReadCandlesticksActivity,
		activity.RegisterOptions{Name: db.ReadCandlesticksActivityName},
	)
	w.RegisterActivityWithOptions(
		a.UpdateCandlesticksActivity,
		activity.RegisterOptions{Name: db.UpdateCandlesticksActivityName},
	)
	w.RegisterActivityWithOptions(
		a.DeleteCandlesticksActivity,
		activity.RegisterOptions{Name: db.DeleteCandlesticksActivityName},
	)
}

// Reset will reset the database.
func (a *Activities) Reset(ctx context.Context) error {
	_, err := a.client.DB.ExecContext(ctx, "DELETE FROM candlesticks")
	if err != nil {
		return fmt.Errorf("deleting candlesticks rows: %w", err)
	}

	return nil
}

// CreateCandlesticksActivity creates the candlesticks.
func (a *Activities) CreateCandlesticksActivity(
	ctx context.Context,
	params db.CreateCandlesticksActivityParams,
) (db.CreateCandlesticksActivityResults, error) {
	// Convert the list of candlesticks from the model to the entity
	listCE, err := entities.FromModelListToEntityList(params.List)
	if err != nil {
		return db.CreateCandlesticksActivityResults{}, err
	}

	// Bulk insert the candlesticks
	_, err = a.client.DB.NamedExecContext(
		ctx,
		`INSERT INTO candlesticks (exchange, pair, period, time, data)
		VALUES (:exchange, :pair, :period, :time, :data)`,
		entities.FromEntitiesToMap(listCE),
	)
	if err != nil {
		return db.CreateCandlesticksActivityResults{}, fmt.Errorf("bulk inserting candlesticks: %w", err)
	}

	return db.CreateCandlesticksActivityResults{}, nil
}

// ReadCandlesticksActivity reads the candlesticks.
func (a *Activities) ReadCandlesticksActivity(
	ctx context.Context,
	params db.ReadCandlesticksActivityParams,
) (db.ReadCandlesticksActivityResults, error) {
	// Set artificial limit if there is none
	if params.Limit == 0 {
		params.Limit = math.MaxInt32
	}

	// Query the candlesticks
	rows, err := a.client.DB.QueryxContext(
		ctx,
		`SELECT exchange, pair, period, time, data
		FROM candlesticks
		WHERE exchange = $1 AND pair = $2 AND period = $3 AND time >= $4 AND time <= $5
		ORDER BY time ASC
		LIMIT $6`,
		params.Exchange,
		params.Pair,
		params.Period,
		params.Start.UTC(),
		params.End.UTC(),
		params.Limit,
	)
	if err != nil {
		return db.ReadCandlesticksActivityResults{}, fmt.Errorf("querying candlesticks: %w", err)
	}
	defer rows.Close()

	// Loop through the results
	list := candlestick.NewList(params.Exchange, params.Pair, params.Period)
	for rows.Next() {
		ce := entities.Candlestick{}
		if err := rows.StructScan(&ce); err != nil {
			return db.ReadCandlesticksActivityResults{}, fmt.Errorf("scanning candlestick: %w", err)
		}

		_, _, _, m, err := ce.ToModel()
		if err != nil {
			return db.ReadCandlesticksActivityResults{}, fmt.Errorf("from candlestick entity to model: %w", err)
		}

		if err := list.Set(m); err != nil {
			return db.ReadCandlesticksActivityResults{}, fmt.Errorf("setting candlestick: %w", err)
		}
	}

	return db.ReadCandlesticksActivityResults{
		List: list,
	}, nil
}

// UpdateCandlesticksActivity updates the candlesticks.
func (a *Activities) UpdateCandlesticksActivity(
	ctx context.Context,
	params db.UpdateCandlesticksActivityParams,
) (db.UpdateCandlesticksActivityResults, error) {
	// Convert the list of candlesticks from the model to the entity
	listCE, err := entities.FromModelListToEntityList(params.List)
	if err != nil {
		return db.UpdateCandlesticksActivityResults{}, err
	}

	// Update the candlesticks
	for _, ce := range listCE {
		res, err := a.client.DB.NamedExecContext(
			ctx,
			`UPDATE candlesticks
			SET data = :data
			WHERE exchange = :exchange AND pair = :pair AND period = :period AND time = :time`,
			ce,
		)
		if err != nil {
			return db.UpdateCandlesticksActivityResults{}, fmt.Errorf("updating candlesticks: %w", err)
		}

		// Check if the candlestick was updated
		affected, err := res.RowsAffected()
		if err != nil {
			return db.UpdateCandlesticksActivityResults{}, fmt.Errorf("getting rows affected: %w", err)
		} else if affected == 0 {
			return db.UpdateCandlesticksActivityResults{}, db.ErrNotFound
		}
	}

	return db.UpdateCandlesticksActivityResults{}, nil
}

// DeleteCandlesticksActivity deletes the candlesticks.
func (a *Activities) DeleteCandlesticksActivity(
	ctx context.Context,
	params db.DeleteCandlesticksActivityParams,
) (db.DeleteCandlesticksActivityResults, error) {
	// Get the times
	times := make([]time.Time, 0, params.List.Data.Len())
	_ = params.List.Loop(func(cs candlestick.Candlestick) (bool, error) {
		times = append(times, cs.Time)
		return false, nil
	})

	query, args, err := sqlx.In(
		`DELETE FROM candlesticks
		WHERE exchange = ? AND pair = ? AND period = ? AND time IN (?)`,
		params.List.Metadata.Exchange,
		params.List.Metadata.Pair,
		params.List.Metadata.Period.String(),
		times,
	)
	if err != nil {
		return db.DeleteCandlesticksActivityResults{}, fmt.Errorf("building query: %w", err)
	}

	query = a.client.DB.Rebind(query)
	_, err = a.client.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return db.DeleteCandlesticksActivityResults{}, fmt.Errorf("deleting candlesticks: %w", err)
	}

	return db.DeleteCandlesticksActivityResults{}, nil
}
