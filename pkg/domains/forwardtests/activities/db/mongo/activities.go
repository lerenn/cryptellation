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
	w.RegisterActivityWithOptions(a.CreateForwardTestActivity,
		activity.RegisterOptions{Name: db.CreateForwardTestActivityName})
	w.RegisterActivityWithOptions(a.ReadForwardTestActivity,
		activity.RegisterOptions{Name: db.ReadForwardTestActivityName})
	w.RegisterActivityWithOptions(a.ListForwardTestsActivity,
		activity.RegisterOptions{Name: db.ListForwardTestsActivityName})
	w.RegisterActivityWithOptions(a.UpdateForwardTestActivity,
		activity.RegisterOptions{Name: db.UpdateForwardTestActivityName})
	w.RegisterActivityWithOptions(a.DeleteForwardTestActivity,
		activity.RegisterOptions{Name: db.DeleteForwardTestActivityName})
}

const (
	// CollectionName is the name of the collection in the database.
	CollectionName = "forwardtests"
)

// CreateForwardTestActivity creates a new forwardtest in the database.
func (a *Activities) CreateForwardTestActivity(
	ctx context.Context,
	params db.CreateForwardTestActivityParams,
) (db.CreateForwardTestActivityResult, error) {
	// Check ID is not nil
	if params.ForwardTest.ID == uuid.Nil {
		return db.CreateForwardTestActivityResult{}, db.ErrNilID
	}

	entity := entities.FromForwardTestModel(params.ForwardTest)

	_, err := a.client.Collection(CollectionName).InsertOne(ctx, entity)
	return db.CreateForwardTestActivityResult{}, err
}

// ReadForwardTestActivity reads a forwardtest from the database.
func (a *Activities) ReadForwardTestActivity(
	ctx context.Context,
	params db.ReadForwardTestActivityParams,
) (db.ReadForwardTestActivityResult, error) {
	var entity entities.ForwardTest

	// Check ID is not nil
	if params.ID == uuid.Nil {
		return db.ReadForwardTestActivityResult{}, db.ErrNilID
	}

	// Get object from database
	err := a.client.
		Collection(CollectionName).
		FindOne(ctx, map[string]any{"_id": params.ID.String()}).
		Decode(&entity)
	if err != nil {
		return db.ReadForwardTestActivityResult{}, err
	}

	// Transform entity to model
	ft, err := entity.ToModel()
	return db.ReadForwardTestActivityResult{
		ForwardTest: ft,
	}, err
}

// ListForwardTestsActivity lists all forwardtests from the database.
func (a *Activities) ListForwardTestsActivity(
	ctx context.Context,
	_ db.ListForwardTestsActivityParams,
) (db.ListForwardTestsActivityResult, error) {
	var models []forwardtest.ForwardTest

	findOptions := options.Find()
	// Sort by `price` field descending
	findOptions.SetSort(bson.D{{Key: "updated_at", Value: -1}})

	// Get objects from database
	cursor, err := a.client.Collection(CollectionName).
		Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return db.ListForwardTestsActivityResult{}, err
	}
	defer cursor.Close(ctx)

	// Transform entities to models
	for cursor.Next(ctx) {
		var entity entities.ForwardTest
		err := cursor.Decode(&entity)
		if err != nil {
			return db.ListForwardTestsActivityResult{}, err
		}

		model, err := entity.ToModel()
		if err != nil {
			return db.ListForwardTestsActivityResult{}, err
		}

		models = append(models, model)
	}

	return db.ListForwardTestsActivityResult{
		ForwardTests: models,
	}, nil
}

// UpdateForwardTestActivity updates a forwardtest in the database.
func (a *Activities) UpdateForwardTestActivity(
	ctx context.Context,
	params db.UpdateForwardTestActivityParams,
) (db.UpdateForwardTestActivityResult, error) {
	// Check ID is not nil
	if params.ForwardTest.ID == uuid.Nil {
		return db.UpdateForwardTestActivityResult{}, db.ErrNilID
	}

	// Update backtest
	entity := entities.FromForwardTestModel(params.ForwardTest)
	_, err := a.client.
		Collection(CollectionName).
		ReplaceOne(ctx, map[string]any{"_id": params.ForwardTest.ID.String()}, entity)
	return db.UpdateForwardTestActivityResult{}, err
}

// DeleteForwardTestActivity deletes a forwardtest from the database.
func (a *Activities) DeleteForwardTestActivity(
	ctx context.Context,
	params db.DeleteForwardTestActivityParams,
) (db.DeleteForwardTestActivityResult, error) {
	// Check ID is not nil
	if params.ID == uuid.Nil {
		return db.DeleteForwardTestActivityResult{}, db.ErrNilID
	}

	// Delete backtest
	_, err := a.client.Collection(CollectionName).
		DeleteOne(ctx, map[string]any{"_id": params.ID.String()})
	return db.DeleteForwardTestActivityResult{}, err
}
