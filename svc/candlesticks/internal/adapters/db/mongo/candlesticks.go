package mongo

import (
	"context"
	"time"

	mongoutil "github.com/lerenn/cryptellation/internal/adapters/db/mongo"

	"github.com/lerenn/cryptellation/svc/candlesticks/internal/adapters/db/mongo/entities"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// CollectionName is the name of the collection in the database
	CollectionName = "candlesticks"
)

func (a *Adapter) CreateCandlesticks(ctx context.Context, cs *candlestick.List) error {
	listCE, err := entities.FromModelListToEntityList(cs)
	if err != nil {
		return err
	}

	_, err = a.client.
		Collection(CollectionName).
		InsertMany(ctx, mongoutil.ToInterfaceList(listCE))
	return err
}

func (a *Adapter) ReadCandlesticks(ctx context.Context, cs *candlestick.List, start, end time.Time, limit uint) error {
	opts := make([]*options.FindOptions, 0, 1)
	if limit != 0 {
		opts = append(opts, options.Find().SetLimit(int64(limit)))
	}

	// Find the candlesticks
	cursor, err := a.client.
		Collection(CollectionName).
		Find(ctx, bson.M{
			"exchange": cs.Exchange,
			"pair":     cs.Pair,
			"period":   cs.Period.String(),
			"time":     bson.M{"$gte": start, "$lte": end},
		}, opts...)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	// Loop through the results
	for cursor.Next(ctx) {
		ce := entities.Candlestick{}
		if err := cursor.Decode(&ce); err != nil {
			return err
		}

		_, _, _, m := ce.ToModel()
		if err := cs.Set(m); err != nil {
			return err
		}
	}

	return nil
}

func (a *Adapter) UpdateCandlesticks(ctx context.Context, cs *candlestick.List) error {
	listCE, err := entities.FromModelListToEntityList(cs)
	if err != nil {
		return err
	}

	for _, ce := range listCE {
		res, err := a.client.
			Collection(CollectionName).
			UpdateOne(ctx, bson.M{
				"exchange": ce.Exchange,
				"pair":     ce.Pair,
				"period":   ce.Period,
				"time":     ce.Time,
			}, bson.M{
				"$set": ce,
			})
		if err != nil {
			return err
		}

		if res.ModifiedCount == 0 {
			return ErrNoDocument
		}
	}

	return nil
}

func (a *Adapter) DeleteCandlesticks(ctx context.Context, cs *candlestick.List) error {
	// Get the times
	times := make([]time.Time, 0, cs.Len())
	_ = cs.Loop(func(cs candlestick.Candlestick) (bool, error) {
		times = append(times, cs.Time)
		return false, nil
	})

	// Delete the candlesticks
	_, err := a.client.
		Collection(CollectionName).
		DeleteMany(ctx, bson.M{
			"exchange": cs.Exchange,
			"pair":     cs.Pair,
			"period":   cs.Period.String(),
			"time":     bson.M{"$in": times},
		})

	return err
}
