package sql

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/v1/pkg/activities"
	"github.com/lerenn/cryptellation/v1/pkg/config"
	"github.com/lerenn/cryptellation/v1/pkg/domains/forwardtests/activities/db"
	"github.com/lerenn/cryptellation/v1/pkg/domains/forwardtests/activities/db/sql/entities"
	"github.com/lerenn/cryptellation/v1/pkg/models/forwardtest"
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

// Register registers the activities to the worker.
func (a *Activities) Register(w worker.Worker) {
	w.RegisterActivityWithOptions(a.CreateForwardtestActivity,
		activity.RegisterOptions{Name: db.CreateForwardtestActivityName})
	w.RegisterActivityWithOptions(a.ReadForwardtestActivity,
		activity.RegisterOptions{Name: db.ReadForwardtestActivityName})
	w.RegisterActivityWithOptions(a.ListForwardtestsActivity,
		activity.RegisterOptions{Name: db.ListForwardtestsActivityName})
	w.RegisterActivityWithOptions(a.UpdateForwardtestActivity,
		activity.RegisterOptions{Name: db.UpdateForwardtestActivityName})
	w.RegisterActivityWithOptions(a.DeleteForwardtestActivity,
		activity.RegisterOptions{Name: db.DeleteForwardtestActivityName})
}

// Reset will reset the database.
func (a *Activities) Reset(ctx context.Context) error {
	_, err := a.client.DB.ExecContext(ctx, "DELETE FROM forwardtests")
	if err != nil {
		return fmt.Errorf("deleting forwardtests rows: %w", err)
	}

	return nil
}

// CreateForwardtestActivity creates a new forwardtest in the database.
func (a *Activities) CreateForwardtestActivity(
	ctx context.Context,
	params db.CreateForwardtestActivityParams,
) (db.CreateForwardtestActivityResult, error) {
	// Check ID is not nil
	if params.Forwardtest.ID == uuid.Nil {
		return db.CreateForwardtestActivityResult{}, db.ErrNilID
	}

	entity, err := entities.FromForwardtestModel(params.Forwardtest)
	if err != nil {
		return db.CreateForwardtestActivityResult{}, fmt.Errorf("converting forwardtest model to entity: %w", err)
	}

	_, err = a.client.DB.NamedExecContext(ctx, `
		INSERT INTO forwardtests (id, updated_at, data)
		VALUES (:id, :updated_at, :data)
	`, entity)
	if err != nil {
		return db.CreateForwardtestActivityResult{}, fmt.Errorf("inserting forwardtest row: %w", err)
	}

	return db.CreateForwardtestActivityResult{}, nil
}

// ReadForwardtestActivity reads a forwardtest from the database.
func (a *Activities) ReadForwardtestActivity(
	ctx context.Context,
	params db.ReadForwardtestActivityParams,
) (db.ReadForwardtestActivityResult, error) {
	var entity entities.Forwardtest

	// Check ID is not nil
	if params.ID == uuid.Nil {
		return db.ReadForwardtestActivityResult{}, db.ErrNilID
	}

	err := a.client.DB.GetContext(ctx, &entity, "SELECT * FROM forwardtests WHERE id = $1", params.ID)
	if err != nil {
		return db.ReadForwardtestActivityResult{}, fmt.Errorf("reading forwardtest: %w", err)
	}

	forwardtest, err := entity.ToModel()
	if err != nil {
		return db.ReadForwardtestActivityResult{}, fmt.Errorf("converting forwardtest entity to model: %w", err)
	}

	return db.ReadForwardtestActivityResult{
		Forwardtest: forwardtest,
	}, nil
}

// ListForwardtestsActivity lists all forwardtests from the database.
func (a *Activities) ListForwardtestsActivity(
	ctx context.Context,
	_ db.ListForwardtestsActivityParams,
) (db.ListForwardtestsActivityResult, error) {
	var entities []entities.Forwardtest

	err := a.client.DB.SelectContext(ctx, &entities, `
		SELECT *
		FROM forwardtests
		ORDER BY updated_at DESC
	`)
	if err != nil {
		return db.ListForwardtestsActivityResult{}, fmt.Errorf("querying forwardtests rows: %w", err)
	}

	models := make([]forwardtest.Forwardtest, 0, len(entities))
	for _, entity := range entities {
		model, err := entity.ToModel()
		if err != nil {
			return db.ListForwardtestsActivityResult{}, fmt.Errorf("converting forwardtest entity to model: %w", err)
		}
		models = append(models, model)
	}

	return db.ListForwardtestsActivityResult{
		Forwardtests: models,
	}, nil
}

// UpdateForwardtestActivity updates a forwardtest in the database.
func (a *Activities) UpdateForwardtestActivity(
	ctx context.Context,
	params db.UpdateForwardtestActivityParams,
) (db.UpdateForwardtestActivityResult, error) {
	// Check ID is not nil
	if params.Forwardtest.ID == uuid.Nil {
		return db.UpdateForwardtestActivityResult{}, db.ErrNilID
	}

	// Update backtest
	entity, err := entities.FromForwardtestModel(params.Forwardtest)
	if err != nil {
		return db.UpdateForwardtestActivityResult{}, fmt.Errorf("converting forwardtest model to entity: %w", err)
	}

	_, err = a.client.DB.ExecContext(ctx, `
		UPDATE forwardtests
		SET updated_at = $1, data = $2
		WHERE id = $3
	`, entity.UpdatedAt, entity.Data, params.Forwardtest.ID)
	if err != nil {
		return db.UpdateForwardtestActivityResult{}, fmt.Errorf("updating forwardtest row: %w", err)
	}

	return db.UpdateForwardtestActivityResult{}, nil
}

// DeleteForwardtestActivity deletes a forwardtest from the database.
func (a *Activities) DeleteForwardtestActivity(
	ctx context.Context,
	params db.DeleteForwardtestActivityParams,
) (db.DeleteForwardtestActivityResult, error) {
	// Check ID is not nil
	if params.ID == uuid.Nil {
		return db.DeleteForwardtestActivityResult{}, db.ErrNilID
	}

	_, err := a.client.DB.ExecContext(ctx, `
		DELETE FROM forwardtests
		WHERE id = $1
	`, params.ID)
	if err != nil {
		return db.DeleteForwardtestActivityResult{}, fmt.Errorf("deleting forwardtest row: %w", err)
	}

	return db.DeleteForwardtestActivityResult{}, nil
}
