package mongo

import (
	"context"

	"github.com/lerenn/cryptellation/v1/internal/activities"
	"github.com/lerenn/cryptellation/v1/internal/domains/candlesticks/activities/db"
	"github.com/lerenn/cryptellation/v1/pkg/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
)

var _ db.Interface = (*Activities)(nil)

// Activities is the activities.
type Activities struct {
	client activities.Mongo
}

const (
	collectionName = "candlesticks"
)

// New creates a new activities.
func New(ctx context.Context, c config.Mongo) (*Activities, error) {
	// Create embedded database access
	db, err := activities.NewMongo(ctx, c)
	if err != nil {
		return nil, err
	}

	// Create a structure
	a := &Activities{
		client: db,
	}

	// Create indexes
	return a, a.CreateIndexes(ctx)
}

// Register registers the activities.
func (a Activities) Register(w worker.Worker) {
	w.RegisterActivityWithOptions(a.CreateCandlesticks, activity.RegisterOptions{Name: db.CreateCandlesticksActivityName})
	w.RegisterActivityWithOptions(a.ReadCandlesticks, activity.RegisterOptions{Name: db.ReadCandlesticksActivityName})
	w.RegisterActivityWithOptions(a.UpdateCandlesticks, activity.RegisterOptions{Name: db.UpdateCandlesticksActivityName})
	w.RegisterActivityWithOptions(a.DeleteCandlesticks, activity.RegisterOptions{Name: db.DeleteCandlesticksActivityName})
}

// CreateIndexes creates the indexes in the database.
func (a *Activities) CreateIndexes(ctx context.Context) error {
	_, err := a.client.
		Collection(collectionName).
		Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "exchange", Value: 1},
			{Key: "pair", Value: 1},
			{Key: "period", Value: 1},
			{Key: "time", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})

	return err
}

// Reset will reset the database.
func (a *Activities) Reset(ctx context.Context) error {
	if err := a.client.DropDatabase(ctx); err != nil {
		return err
	}

	return a.CreateIndexes(ctx)
}
