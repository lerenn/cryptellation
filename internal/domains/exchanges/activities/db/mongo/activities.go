package mongo

import (
	"context"

	"github.com/lerenn/cryptellation/v1/internal/activities"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"

	"github.com/lerenn/cryptellation/v1/pkg/config"

	"github.com/lerenn/cryptellation/v1/internal/domains/exchanges/activities/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ db.Interface = (*Activities)(nil)

type Activities struct {
	client activities.Mongo
}

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
	return a, a.CreateIndexes()
}

// Register registers the activities.
func (a Activities) Register(w worker.Worker) {
	w.RegisterActivityWithOptions(a.CreateExchanges, activity.RegisterOptions{Name: db.CreateExchangesActivityName})
	w.RegisterActivityWithOptions(a.ReadExchanges, activity.RegisterOptions{Name: db.ReadExchangesActivityName})
	w.RegisterActivityWithOptions(a.UpdateExchanges, activity.RegisterOptions{Name: db.UpdateExchangesActivityName})
	w.RegisterActivityWithOptions(a.DeleteExchanges, activity.RegisterOptions{Name: db.DeleteExchangesActivityName})
}

func (a *Activities) CreateIndexes() error {
	_, err := a.client.
		Collection(CollectionName).
		Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "name", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})

	return err
}

func (a *Activities) Reset(ctx context.Context) error {
	if err := a.client.DropDatabase(ctx); err != nil {
		return err
	}

	return a.CreateIndexes()
}
