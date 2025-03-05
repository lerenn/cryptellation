package mongo

import (
	"context"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/v1/pkg/activities"
	"github.com/lerenn/cryptellation/v1/pkg/config"
	"github.com/lerenn/cryptellation/v1/pkg/domains/forwardtests/activities/db"
	"github.com/lerenn/cryptellation/v1/pkg/domains/forwardtests/activities/db/mongo/entities"
	"github.com/lerenn/cryptellation/v1/pkg/models/forwardtest"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
)

var _ db.DB = (*Activities)(nil)

// Activities is the database access for the forwardtests domain.
type Activities struct {
	client activities.Mongo
}

// New creates a new Activities instance.
func New(ctx context.Context, c config.Mongo) (*Activities, error) {
	// Create embedded database access
	db, err := activities.NewMongo(ctx, c)

	// Return database access
	return &Activities{
		client: db,
	}, err
}

// CreateIndexes creates the indexes for the database.
func (a *Activities) CreateIndexes(ctx context.Context) error {
	index := mongo.IndexModel{Keys: map[string]int{"updated_at": 1}}

	// Create indexes
	_, err := a.client.
		Collection(CollectionName).
		Indexes().
		CreateOne(ctx, index)
	if err != nil {
		return err
	}

	return nil
}

// Reset drops the database and recreates the indexes.
func (a *Activities) Reset(ctx context.Context) error {
	if err := a.client.DropDatabase(ctx); err != nil {
		return err
	}

	return a.CreateIndexes(ctx)
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

const (
	// CollectionName is the name of the collection in the database.
	CollectionName = "forwardtests"
)

// CreateForwardtestActivity creates a new forwardtest in the database.
func (a *Activities) CreateForwardtestActivity(
	ctx context.Context,
	params db.CreateForwardtestActivityParams,
) (db.CreateForwardtestActivityResult, error) {
	// Check ID is not nil
	if params.Forwardtest.ID == uuid.Nil {
		return db.CreateForwardtestActivityResult{}, db.ErrNilID
	}

	entity := entities.FromForwardtestModel(params.Forwardtest)

	_, err := a.client.Collection(CollectionName).InsertOne(ctx, entity)
	return db.CreateForwardtestActivityResult{}, err
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

	// Get object from database
	err := a.client.
		Collection(CollectionName).
		FindOne(ctx, map[string]any{"_id": params.ID.String()}).
		Decode(&entity)
	if err != nil {
		return db.ReadForwardtestActivityResult{}, err
	}

	// Transform entity to model
	ft, err := entity.ToModel()
	return db.ReadForwardtestActivityResult{
		Forwardtest: ft,
	}, err
}

// ListForwardtestsActivity lists all forwardtests from the database.
func (a *Activities) ListForwardtestsActivity(
	ctx context.Context,
	_ db.ListForwardtestsActivityParams,
) (db.ListForwardtestsActivityResult, error) {
	var models []forwardtest.Forwardtest

	findOptions := options.Find()
	// Sort by `price` field descending
	findOptions.SetSort(bson.D{{Key: "updated_at", Value: -1}})

	// Get objects from database
	cursor, err := a.client.Collection(CollectionName).
		Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return db.ListForwardtestsActivityResult{}, err
	}
	defer cursor.Close(ctx)

	// Transform entities to models
	for cursor.Next(ctx) {
		var entity entities.Forwardtest
		err := cursor.Decode(&entity)
		if err != nil {
			return db.ListForwardtestsActivityResult{}, err
		}

		model, err := entity.ToModel()
		if err != nil {
			return db.ListForwardtestsActivityResult{}, err
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
	entity := entities.FromForwardtestModel(params.Forwardtest)
	_, err := a.client.
		Collection(CollectionName).
		ReplaceOne(ctx, map[string]any{"_id": params.Forwardtest.ID.String()}, entity)
	return db.UpdateForwardtestActivityResult{}, err
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

	// Delete backtest
	_, err := a.client.Collection(CollectionName).
		DeleteOne(ctx, map[string]any{"_id": params.ID.String()})
	return db.DeleteForwardtestActivityResult{}, err
}
