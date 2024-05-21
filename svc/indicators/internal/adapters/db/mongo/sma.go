package mongo

import (
	"context"

	"github.com/lerenn/cryptellation/pkg/models/timeserie"
	"github.com/lerenn/cryptellation/svc/indicators/internal/adapters/db/mongo/entities"
	"github.com/lerenn/cryptellation/svc/indicators/internal/app/ports/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// CollectionName is the name of the collection in the database
	CollectionName = "sma"
)

func (a *Adapter) GetSMA(ctx context.Context, payload db.ReadSMAPayload) (*timeserie.TimeSerie[float64], error) {
	// Create a filter
	filter := bson.M{
		"exchange":     payload.Exchange,
		"pair":         payload.Pair,
		"period":       payload.Period,
		"periodNumber": payload.PeriodNumber,
		"priceType":    payload.PriceType,
		"time": bson.M{
			"$gte": payload.Start,
			"$lte": payload.End,
		},
	}

	// Create a cursor
	cursor, err := a.client.Collection(CollectionName).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Create a time serie
	ts := timeserie.New[float64]()
	for cursor.Next(ctx) {
		sma := entities.SimpleMovingAverage{}
		if err := cursor.Decode(&sma); err != nil {
			return nil, err
		}
		_ = ts.Set(sma.Time, sma.Price)
	}

	return ts, nil
}

func (a *Adapter) UpsertSMA(ctx context.Context, payload db.WriteSMAPayload) error {
	// Create entities
	ents := entities.FromModelListToEntityList(
		payload.Exchange,
		payload.Pair,
		payload.Period,
		payload.PeriodNumber,
		payload.PriceType,
		payload.TimeSerie)

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
		_, err := a.client.Collection(CollectionName).UpdateOne(ctx, filter, bson.M{"$set": ent}, options.Update().SetUpsert(true))
		if err != nil {
			return err
		}
	}

	return nil
}
