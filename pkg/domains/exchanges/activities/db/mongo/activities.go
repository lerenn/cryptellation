package mongo

import (
	"context"

	"github.com/lerenn/cryptellation/v1/pkg/activities"
	"github.com/lerenn/cryptellation/v1/pkg/config"
	"github.com/lerenn/cryptellation/v1/pkg/domains/exchanges/activities/db"
	"github.com/lerenn/cryptellation/v1/pkg/domains/exchanges/activities/db/mongo/entities"
	"github.com/lerenn/cryptellation/v1/pkg/models/exchange"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
)

const (
	// CollectionName is the name of the collection in the database.
	CollectionName = "exchanges"
)

var _ db.DB = (*Activities)(nil)

// Activities regroups mongo activities.
type Activities struct {
	client activities.Mongo
}

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
func (a *Activities) Register(w worker.Worker) {
	w.RegisterActivityWithOptions(a.CreateExchangesActivity, activity.RegisterOptions{Name: db.CreateExchangesActivityName})
	w.RegisterActivityWithOptions(a.ReadExchangesActivity, activity.RegisterOptions{Name: db.ReadExchangesActivityName})
	w.RegisterActivityWithOptions(a.UpdateExchangesActivity, activity.RegisterOptions{Name: db.UpdateExchangesActivityName})
	w.RegisterActivityWithOptions(a.DeleteExchangesActivity, activity.RegisterOptions{Name: db.DeleteExchangesActivityName})
}

// CreateIndexes creates the indexes.
func (a *Activities) CreateIndexes(ctx context.Context) error {
	_, err := a.client.
		Collection(CollectionName).
		Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "name", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})

	return err
}

// Reset drops the database and recreates the indexes.
func (a *Activities) Reset(ctx context.Context) error {
	if err := a.client.DropDatabase(ctx); err != nil {
		return err
	}

	return a.CreateIndexes(ctx)
}

// CreateExchangesActivity will create the exchanges in the database.
func (a *Activities) CreateExchangesActivity(
	ctx context.Context,
	params db.CreateExchangesActivityParams,
) (db.CreateExchangesActivityResults, error) {
	ents := make([]interface{}, len(params.Exchanges))
	for i, e := range params.Exchanges {
		ents[i] = entities.ExchangeFromModel(e)
	}

	_, err := a.client.Collection(CollectionName).InsertMany(ctx, ents)
	return db.CreateExchangesActivityResults{}, err
}

// ReadExchangesActivity will read the exchanges from the database.
func (a *Activities) ReadExchangesActivity(
	ctx context.Context,
	params db.ReadExchangesActivityParams,
) (db.ReadExchangesActivityResults, error) {
	filter := bson.M{}
	if len(params.Names) > 0 {
		filter["name"] = bson.M{"$in": params.Names}
	}

	cur, err := a.client.Collection(CollectionName).Find(ctx, filter)
	if err != nil {
		return db.ReadExchangesActivityResults{}, err
	}
	defer cur.Close(ctx)

	var ents []entities.Exchange
	if err := cur.All(ctx, &ents); err != nil {
		return db.ReadExchangesActivityResults{}, err
	}

	exchanges := make([]exchange.Exchange, len(ents))
	for i, e := range ents {
		exchanges[i] = e.ToModel()
	}

	return db.ReadExchangesActivityResults{
		Exchanges: exchanges,
	}, nil
}

// UpdateExchangesActivity will update the exchanges in the database.
func (a *Activities) UpdateExchangesActivity(
	ctx context.Context,
	params db.UpdateExchangesActivityParams,
) (db.UpdateExchangesActivityResults, error) {
	ents := make([]interface{}, len(params.Exchanges))
	for i, e := range params.Exchanges {
		ents[i] = entities.ExchangeFromModel(e)
	}

	for _, e := range params.Exchanges {
		filter := bson.M{"name": e.Name}
		_, err := a.client.Collection(CollectionName).ReplaceOne(ctx, filter, entities.ExchangeFromModel(e))
		if err != nil {
			return db.UpdateExchangesActivityResults{}, err
		}
	}

	return db.UpdateExchangesActivityResults{}, nil
}

// DeleteExchangesActivity will delete the exchanges from the database.
func (a *Activities) DeleteExchangesActivity(
	ctx context.Context,
	params db.DeleteExchangesActivityParams,
) (db.DeleteExchangesActivityResults, error) {
	filter := bson.M{"name": bson.M{"$in": params.Names}}
	_, err := a.client.Collection(CollectionName).DeleteMany(ctx, filter)
	return db.DeleteExchangesActivityResults{}, err
}
