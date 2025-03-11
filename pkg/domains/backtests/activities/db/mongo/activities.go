package mongo

import (
	"context"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/v1/pkg/activities"
	"github.com/lerenn/cryptellation/v1/pkg/config"
	"github.com/lerenn/cryptellation/v1/pkg/domains/backtests/activities/db"
	"github.com/lerenn/cryptellation/v1/pkg/domains/backtests/activities/db/mongo/entities"
	"github.com/lerenn/cryptellation/v1/pkg/models/backtest"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
)

const (
	// CollectionName is the name of the collection in the database.
	CollectionName = "backtests"
)

var _ db.DB = (*Activities)(nil)

// Activities is the database access.
type Activities struct {
	client activities.Mongo
}

// New creates a new database access.
func New(ctx context.Context, c config.Mongo) (*Activities, error) {
	// Create embedded database access
	db, err := activities.NewMongo(ctx, c)

	// Return database access
	return &Activities{
		client: db,
	}, err
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

// CreateIndexes creates the indexes in the database.
func (a *Activities) CreateIndexes() error {
	// No index
	return nil
}

// Reset drops the database and recreates the indexes.
func (a *Activities) Reset(ctx context.Context) error {
	if err := a.client.DropDatabase(ctx); err != nil {
		return err
	}

	return a.CreateIndexes()
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

	// Create backtest
	entity := entities.FromBacktestModel(params.Backtest)
	_, err := a.client.Collection(CollectionName).InsertOne(ctx, entity)
	return db.CreateBacktestActivityResults{}, err
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

	// Get object from database
	err := a.client.
		Collection(CollectionName).
		FindOne(ctx, map[string]any{"_id": params.ID.String()}).
		Decode(&entity)
	if err != nil {
		return db.ReadBacktestActivityResults{}, err
	}

	// Transform entity to model
	bt, err := entity.ToModel()
	if err != nil {
		return db.ReadBacktestActivityResults{}, err
	}

	return db.ReadBacktestActivityResults{Backtest: bt}, nil
}

// ListBacktestsActivity lists backtests from the database.
func (a *Activities) ListBacktestsActivity(
	ctx context.Context,
	_ db.ListBacktestsActivityParams,
) (db.ListBacktestsActivityResults, error) {
	var bts []backtest.Backtest

	// Get objects from database
	cursor, err := a.client.Collection(CollectionName).Find(ctx, map[string]any{})
	if err != nil {
		return db.ListBacktestsActivityResults{}, err
	}
	defer cursor.Close(ctx)

	// Transform entities to models
	for cursor.Next(ctx) {
		var entity entities.Backtest
		if err := cursor.Decode(&entity); err != nil {
			return db.ListBacktestsActivityResults{}, err
		}

		bt, err := entity.ToModel()
		if err != nil {
			return db.ListBacktestsActivityResults{}, err
		}

		bts = append(bts, bt)
	}

	return db.ListBacktestsActivityResults{
		Backtests: bts,
	}, nil
}

// UpdateBacktestActivity updates a backtest in the database.
func (a *Activities) UpdateBacktestActivity(
	ctx context.Context,
	params db.UpdateBacktestActivityParams,
) (db.UpdateBacktestActivityResults, error) {
	// Check ID is not nil
	if params.Backtest.ID == uuid.Nil {
		return db.UpdateBacktestActivityResults{}, db.ErrNilID
	}

	// Update backtest
	entity := entities.FromBacktestModel(params.Backtest)
	_, err := a.client.
		Collection(CollectionName).
		ReplaceOne(ctx, map[string]any{"_id": params.Backtest.ID.String()}, entity)
	return db.UpdateBacktestActivityResults{}, err
}

// DeleteBacktestActivity deletes a backtest from the database.
func (a *Activities) DeleteBacktestActivity(
	ctx context.Context,
	params db.DeleteBacktestActivityParams,
) (db.DeleteBacktestActivityResults, error) {
	// Check ID is not nil
	if params.ID == uuid.Nil {
		return db.DeleteBacktestActivityResults{}, db.ErrNilID
	}

	// Delete backtest
	_, err := a.client.Collection(CollectionName).DeleteOne(ctx, map[string]any{"_id": params.ID.String()})
	return db.DeleteBacktestActivityResults{}, err
}
