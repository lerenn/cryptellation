package mongo

import (
	"context"

	"github.com/lerenn/cryptellation/v1/pkg/activities"
	"github.com/lerenn/cryptellation/v1/pkg/config"
	"github.com/lerenn/cryptellation/v1/pkg/domains/indicators/activities/db"
	"github.com/lerenn/cryptellation/v1/pkg/domains/indicators/activities/db/mongo/entities"
	"github.com/lerenn/cryptellation/v1/pkg/models/timeserie"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
)

var _ db.DB = (*Activities)(nil)

// Activities is the database access for the activities.
type Activities struct {
	client activities.Mongo
}

// New creates a new Activities instance.
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

// CreateIndexes creates the indexes in the database.
func (a *Activities) CreateIndexes(ctx context.Context) error {
	_, err := a.client.
		Collection(CollectionName).
		Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "exchange", Value: 1},
			{Key: "pair", Value: 1},
			{Key: "period", Value: 1},
			{Key: "period_number", Value: 1},
			{Key: "price_type", Value: 1},
			{Key: "time", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})

	return err
}

// Reset resets the database.
func (a *Activities) Reset(ctx context.Context) error {
	if err := a.client.DropDatabase(ctx); err != nil {
		return err
	}

	return a.CreateIndexes(ctx)
}

// Register registers the activities.
func (a *Activities) Register(w worker.Worker) {
	w.RegisterActivityWithOptions(
		a.ReadSMAActivity,
		activity.RegisterOptions{Name: db.ReadSMAActivityName},
	)
	w.RegisterActivityWithOptions(
		a.UpsertSMAActivity,
		activity.RegisterOptions{Name: db.UpsertSMAActivityName},
	)
}

const (
	// CollectionName is the name of the collection in the database.
	CollectionName = "sma"
)

// ReadSMAActivity reads the SMA points from the database.
func (a *Activities) ReadSMAActivity(
	ctx context.Context,
	params db.ReadSMAActivityParams,
) (db.ReadSMAActivityResults, error) {
	// Create a filter
	filter := bson.M{
		"exchange":     params.Exchange,
		"pair":         params.Pair,
		"period":       params.Period,
		"periodNumber": params.PeriodNumber,
		"priceType":    params.PriceType,
		"time": bson.M{
			"$gte": params.Start,
			"$lte": params.End,
		},
	}

	// Create a cursor
	cursor, err := a.client.Collection(CollectionName).Find(ctx, filter)
	if err != nil {
		return db.ReadSMAActivityResults{}, err
	}
	defer cursor.Close(ctx)

	// Create a time serie
	ts := timeserie.New[float64]()
	for cursor.Next(ctx) {
		sma := entities.SimpleMovingAverage{}
		if err := cursor.Decode(&sma); err != nil {
			return db.ReadSMAActivityResults{}, err
		}
		_ = ts.Set(sma.Time, sma.Price)
	}

	return db.ReadSMAActivityResults{
		Data: ts,
	}, nil
}

// UpsertSMAActivity upserts the SMA points in the database.
func (a *Activities) UpsertSMAActivity(
	ctx context.Context,
	params db.UpsertSMAActivityParams,
) (db.UpsertSMAActivityResults, error) {
	// Create entities
	ents := entities.FromModelListToEntityList(
		params.Exchange,
		params.Pair,
		params.Period,
		params.PeriodNumber,
		params.PriceType,
		params.TimeSerie)

	for _, ent := range ents {
		// Create a filter
		filter := bson.M{
			"exchange":     ent.Exchange,
			"pair":         ent.Pair,
			"period":       ent.Period,
			"periodNumber": ent.PeriodNumber,
			"priceType":    ent.PriceType,
			"time":         ent.Time,
		}

		// Upsert the document
		_, err := a.client.Collection(CollectionName).
			UpdateOne(ctx, filter, bson.M{
				"$set": ent,
			}, options.Update().SetUpsert(true))
		if err != nil {
			return db.UpsertSMAActivityResults{}, err
		}
	}

	return db.UpsertSMAActivityResults{}, nil
}
