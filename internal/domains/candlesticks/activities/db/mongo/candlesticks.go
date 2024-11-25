package mongo

import (
	"context"
	"time"

	"github.com/lerenn/cryptellation/v1/internal/domains/candlesticks/activities/db"
	"github.com/lerenn/cryptellation/v1/internal/domains/candlesticks/activities/db/mongo/entities"
	"github.com/lerenn/cryptellation/v1/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/v1/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateCandlesticks creates the candlesticks.
func (a *Activities) CreateCandlesticks(
	ctx context.Context,
	params db.CreateCandlesticksParams,
) (db.CreateCandlesticksResult, error) {
	listCE, err := entities.FromModelListToEntityList(params.List)
	if err != nil {
		return db.CreateCandlesticksResult{}, err
	}

	_, err = a.client.
		Collection(collectionName).
		InsertMany(ctx, utils.ToInterfaceList(listCE))
	return db.CreateCandlesticksResult{}, err
}

// ReadCandlesticks reads the candlesticks.
func (a *Activities) ReadCandlesticks(
	ctx context.Context,
	params db.ReadCandlesticksParams,
) (db.ReadCandlesticksResult, error) {
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
		return db.ReadCandlesticksResult{}, err
	}
	defer cursor.Close(ctx)

	// Loop through the results
	list := candlestick.NewList(params.Exchange, params.Pair, params.Period)
	for cursor.Next(ctx) {
		ce := entities.Candlestick{}
		if err := cursor.Decode(&ce); err != nil {
			return db.ReadCandlesticksResult{}, err
		}

		_, _, _, m := ce.ToModel()
		if err := list.Set(m); err != nil {
			return db.ReadCandlesticksResult{}, err
		}
	}

	return db.ReadCandlesticksResult{
		List: list,
	}, nil
}

// UpdateCandlesticks updates the candlesticks.
func (a *Activities) UpdateCandlesticks(
	ctx context.Context,
	params db.UpdateCandlesticksParams,
) (db.UpdateCandlesticksResult, error) {
	listCE, err := entities.FromModelListToEntityList(params.List)
	if err != nil {
		return db.UpdateCandlesticksResult{}, err
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
			return db.UpdateCandlesticksResult{}, err
		}

		if res.ModifiedCount == 0 {
			return db.UpdateCandlesticksResult{}, ErrNoDocument
		}
	}

	return db.UpdateCandlesticksResult{}, nil
}

// DeleteCandlesticks deletes the candlesticks.
func (a *Activities) DeleteCandlesticks(
	ctx context.Context,
	params db.DeleteCandlesticksParams,
) (db.DeleteCandlesticksResult, error) {
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

	return db.DeleteCandlesticksResult{}, err
}
