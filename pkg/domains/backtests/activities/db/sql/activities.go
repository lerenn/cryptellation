package sql

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/v1/pkg/activities"
	"github.com/lerenn/cryptellation/v1/pkg/config"
	"github.com/lerenn/cryptellation/v1/pkg/domains/backtests/activities/db"
	"github.com/lerenn/cryptellation/v1/pkg/domains/backtests/activities/db/sql/entities"
	"github.com/lerenn/cryptellation/v1/pkg/models/backtest"
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

// Register registers the activities to the worker.
func (a *Activities) Register(w worker.Worker) {
	w.RegisterActivityWithOptions(
		a.CreateBacktestActivity,
		activity.RegisterOptions{Name: db.CreateBacktestActivityName},
	)

	w.RegisterActivityWithOptions(
		a.ReadBacktestActivity,
		activity.RegisterOptions{Name: db.ReadBacktestActivityName},
	)

	w.RegisterActivityWithOptions(
		a.ListBacktestsActivity,
		activity.RegisterOptions{Name: db.ListBacktestsActivityName},
	)

	w.RegisterActivityWithOptions(
		a.UpdateBacktestActivity,
		activity.RegisterOptions{Name: db.UpdateBacktestActivityName},
	)

	w.RegisterActivityWithOptions(
		a.DeleteBacktestActivity,
		activity.RegisterOptions{Name: db.DeleteBacktestActivityName},
	)
}

// Reset will reset the database.
func (a *Activities) Reset(ctx context.Context) error {
	_, err := a.client.DB.ExecContext(ctx, "DELETE FROM backtests")
	if err != nil {
		return fmt.Errorf("deleting backtests rows: %w", err)
	}

	return nil
}

// CreateBacktestActivity creates a backtest in the database.
func (a *Activities) CreateBacktestActivity(
	ctx context.Context,
	params db.CreateBacktestActivityParams,
) (db.CreateBacktestActivityResults, error) {
	// Check ID is not nil
	if params.Backtest.ID == uuid.Nil {
		return db.CreateBacktestActivityResults{}, db.ErrNilID
	}

	// Change backtest model to entity
	entity, err := entities.FromBacktestModel(params.Backtest)
	if err != nil {
		return db.CreateBacktestActivityResults{}, err
	}

	// Insert the backtest
	_, err = a.client.DB.NamedExecContext(
		ctx,
		`INSERT INTO backtests (id, data)
		VALUES (:id, :data)`,
		entity)
	if err != nil {
		return db.CreateBacktestActivityResults{}, fmt.Errorf("inserting backtest: %w", err)
	}

	return db.CreateBacktestActivityResults{}, nil
}

// ReadBacktestActivity reads a backtest from the database.
func (a *Activities) ReadBacktestActivity(
	ctx context.Context,
	params db.ReadBacktestActivityParams,
) (db.ReadBacktestActivityResults, error) {
	var entity entities.Backtest

	// Check ID is not nil
	if params.ID == uuid.Nil {
		return db.ReadBacktestActivityResults{}, db.ErrNilID
	}

	// Read the backtest
	err := a.client.DB.GetContext(ctx, &entity, "SELECT * FROM backtests WHERE id = $1", params.ID)
	if err != nil {
		return db.ReadBacktestActivityResults{}, fmt.Errorf("reading backtest: %w", err)
	}

	// Convert the entity to the model
	m, err := entity.ToModel()
	if err != nil {
		return db.ReadBacktestActivityResults{}, fmt.Errorf("converting entity to model: %w", err)
	}

	return db.ReadBacktestActivityResults{Backtest: m}, nil
}

// ListBacktestsActivity lists backtests from the database.
func (a *Activities) ListBacktestsActivity(
	ctx context.Context,
	_ db.ListBacktestsActivityParams,
) (db.ListBacktestsActivityResults, error) {
	var entities []entities.Backtest

	// Read the backtests
	err := a.client.DB.SelectContext(ctx, &entities, "SELECT * FROM backtests")
	if err != nil {
		return db.ListBacktestsActivityResults{}, fmt.Errorf("reading backtests: %w", err)
	}

	// Convert the entities to models
	models := make([]backtest.Backtest, 0, len(entities))
	for _, e := range entities {
		m, err := e.ToModel()
		if err != nil {
			return db.ListBacktestsActivityResults{}, fmt.Errorf("converting entity to model: %w", err)
		}
		models = append(models, m)
	}

	return db.ListBacktestsActivityResults{Backtests: models}, nil
}

// UpdateBacktestActivity updates the backtest in the database.
func (a *Activities) UpdateBacktestActivity(
	ctx context.Context,
	params db.UpdateBacktestActivityParams,
) (db.UpdateBacktestActivityResults, error) {
	// Check ID is not nil
	if params.Backtest.ID == uuid.Nil {
		return db.UpdateBacktestActivityResults{}, db.ErrNilID
	}

	// Change the backtest model to entity
	entity, err := entities.FromBacktestModel(params.Backtest)
	if err != nil {
		return db.UpdateBacktestActivityResults{}, err
	}

	// Update the backtest
	_, err = a.client.DB.NamedExecContext(
		ctx,
		`UPDATE backtests
		SET data = :data
		WHERE id = :id`,
		entity)
	if err != nil {
		return db.UpdateBacktestActivityResults{}, fmt.Errorf("updating backtest: %w", err)
	}

	return db.UpdateBacktestActivityResults{}, nil
}

// DeleteBacktestActivity deletes the backtest from the database.
func (a *Activities) DeleteBacktestActivity(
	ctx context.Context,
	params db.DeleteBacktestActivityParams,
) (db.DeleteBacktestActivityResults, error) {
	// Check ID is not nil
	if params.ID == uuid.Nil {
		return db.DeleteBacktestActivityResults{}, db.ErrNilID
	}

	// Delete the backtest
	_, err := a.client.DB.ExecContext(ctx, "DELETE FROM backtests WHERE id = $1", params.ID)
	if err != nil {
		return db.DeleteBacktestActivityResults{}, fmt.Errorf("deleting backtest: %w", err)
	}

	return db.DeleteBacktestActivityResults{}, nil
}
