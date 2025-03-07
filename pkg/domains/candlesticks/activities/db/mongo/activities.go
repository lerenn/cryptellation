package mongo

import (
	"context"
	"time"

	"github.com/lerenn/cryptellation/v1/pkg/activities"
	"github.com/lerenn/cryptellation/v1/pkg/config"
	"github.com/lerenn/cryptellation/v1/pkg/domains/candlesticks/activities/db"
	"github.com/lerenn/cryptellation/v1/pkg/domains/candlesticks/activities/db/mongo/entities"
	"github.com/lerenn/cryptellation/v1/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/v1/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
)

var _ db.DB = (*Activities)(nil)

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

// CreateCandlesticksActivity creates the candlesticks.
func (a *Activities) CreateCandlesticksActivity(
	ctx context.Context,
	params db.CreateCandlesticksActivityParams,
) (db.CreateCandlesticksActivityResults, error) {
	listCE, err := entities.FromModelListToEntityList(params.List)
	if err != nil {
		return db.CreateCandlesticksActivityResults{}, err
	}

	_, err = a.client.
		Collection(collectionName).
		InsertMany(ctx, utils.ToInterfaceList(listCE))
	return db.CreateCandlesticksActivityResults{}, err
}

// ReadCandlesticksActivity reads the candlesticks.
func (a *Activities) ReadCandlesticksActivity(
	ctx context.Context,
	params db.ReadCandlesticksActivityParams,
) (db.ReadCandlesticksActivityResults, error) {
	opts := make([]*options.FindOptions, 0, 1)
	if params.Limit != 0 {
		opts = append(opts, options.Find().SetLimit(int64(params.Limit)))
	}

	// Find the candlesticks
	cursor, err := a.client.
		Collection(collectionName).
		Find(ctx, bson.M{
			"exchange": params.Exchange,
			"pair":     params.Pair,
			"period":   params.Period.String(),
			"time":     bson.M{"$gte": params.Start, "$lte": params.End},
		}, opts...)
	if err != nil {
		return db.ReadCandlesticksActivityResults{}, err
	}
	defer cursor.Close(ctx)

	// Loop through the results
	list := candlestick.NewList(params.Exchange, params.Pair, params.Period)
	for cursor.Next(ctx) {
		ce := entities.Candlestick{}
		if err := cursor.Decode(&ce); err != nil {
			return db.ReadCandlesticksActivityResults{}, err
		}

		_, _, _, m := ce.ToModel()
		if err := list.Set(m); err != nil {
			return db.ReadCandlesticksActivityResults{}, err
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
	listCE, err := entities.FromModelListToEntityList(params.List)
	if err != nil {
		return db.UpdateCandlesticksActivityResults{}, err
	}

	for _, ce := range listCE {
		res, err := a.client.
			Collection(collectionName).
			UpdateOne(ctx, bson.M{
				"exchange": ce.Exchange,
				"pair":     ce.Pair,
				"period":   ce.Period,
				"time":     ce.Time,
			}, bson.M{
				"$set": ce,
			})
		if err != nil {
			return db.UpdateCandlesticksActivityResults{}, err
		}

		if res.ModifiedCount == 0 {
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

	// Delete the candlesticks
	_, err := a.client.
		Collection(collectionName).
		DeleteMany(ctx, bson.M{
			"exchange": params.List.Metadata.Exchange,
			"pair":     params.List.Metadata.Pair,
			"period":   params.List.Metadata.Period.String(),
			"time":     bson.M{"$in": times},
		})

	return db.DeleteCandlesticksActivityResults{}, err
}
