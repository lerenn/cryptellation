package mongo

import (
	"context"

	"github.com/lerenn/cryptellation/svc/forwardtests/internal/adapters/db/mongo/entities"
	"github.com/lerenn/cryptellation/svc/forwardtests/internal/app/ports/db"
	"github.com/lerenn/cryptellation/svc/forwardtests/pkg/forwardtest"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// CollectionName is the name of the collection in the database
	CollectionName = "forwardtests"
)

func (mongo *Adapter) CreateForwardTest(ctx context.Context, ft forwardtest.ForwardTest) error {
	// Check ID is not nil
	if ft.ID == uuid.Nil {
		return ErrNilID
	}

	entity := entities.FromForwardTestModel(ft)

	_, err := mongo.client.Collection(CollectionName).InsertOne(ctx, entity)
	return err
}

func (mongo *Adapter) ReadForwardTest(ctx context.Context, id uuid.UUID) (forwardtest.ForwardTest, error) {
	var entity entities.ForwardTest

	// Check ID is not nil
	if id == uuid.Nil {
		return forwardtest.ForwardTest{}, ErrNilID
	}

	// Get object from database
	err := mongo.client.
		Collection(CollectionName).
		FindOne(ctx, map[string]any{"_id": id.String()}).
		Decode(&entity)
	if err != nil {
		return forwardtest.ForwardTest{}, err
	}

	// Transform entity to model
	return entity.ToModel()
}

func (mongo *Adapter) ListForwardTests(ctx context.Context, filters db.ListFilters) ([]forwardtest.ForwardTest, error) {
	var models []forwardtest.ForwardTest

	findOptions := options.Find()
	// Sort by `price` field descending
	findOptions.SetSort(bson.D{{Key: "updated_at", Value: -1}})

	// Get objects from database
	cursor, err := mongo.client.Collection(CollectionName).Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Transform entities to models
	for cursor.Next(ctx) {
		var entity entities.ForwardTest
		err := cursor.Decode(&entity)
		if err != nil {
			return nil, err
		}

		model, err := entity.ToModel()
		if err != nil {
			return nil, err
		}

		models = append(models, model)
	}

	return models, nil
}

func (mongo *Adapter) UpdateForwardTest(ctx context.Context, ft forwardtest.ForwardTest) error {
	// Check ID is not nil
	if ft.ID == uuid.Nil {
		return ErrNilID
	}

	// Update backtest
	entity := entities.FromForwardTestModel(ft)
	_, err := mongo.client.
		Collection(CollectionName).
		ReplaceOne(ctx, map[string]any{"_id": ft.ID.String()}, entity)
	return err
}

func (mongo *Adapter) DeleteForwardTest(ctx context.Context, id uuid.UUID) error {
	// Check ID is not nil
	if id == uuid.Nil {
		return ErrNilID
	}

	// Delete backtest
	_, err := mongo.client.Collection(CollectionName).DeleteOne(ctx, map[string]any{"_id": id.String()})
	return err
}
