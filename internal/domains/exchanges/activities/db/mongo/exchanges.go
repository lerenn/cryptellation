package mongo

import (
	"context"

	"github.com/lerenn/cryptellation/v1/internal/domains/exchanges/activities/db"
	"github.com/lerenn/cryptellation/v1/internal/domains/exchanges/activities/db/mongo/entities"
	"github.com/lerenn/cryptellation/v1/pkg/models/exchange"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	// CollectionName is the name of the collection in the database.
	CollectionName = "exchanges"
)

// CreateExchanges will create the exchanges in the database.
func (a *Activities) CreateExchanges(
	ctx context.Context,
	params db.CreateExchangesParams,
) (db.CreateExchangesResult, error) {
	ents := make([]interface{}, len(params.Exchanges))
	for i, e := range params.Exchanges {
		ents[i] = entities.ExchangeFromModel(e)
	}

	_, err := a.client.Collection(CollectionName).InsertMany(ctx, ents)
	return db.CreateExchangesResult{}, err
}

// ReadExchanges will read the exchanges from the database.
func (a *Activities) ReadExchanges(
	ctx context.Context,
	params db.ReadExchangesParams,
) (db.ReadExchangesResult, error) {
	filter := bson.M{}
	if len(params.Names) > 0 {
		filter["name"] = bson.M{"$in": params.Names}
	}

	cur, err := a.client.Collection(CollectionName).Find(ctx, filter)
	if err != nil {
		return db.ReadExchangesResult{}, err
	}
	defer cur.Close(ctx)

	var ents []entities.Exchange
	if err := cur.All(ctx, &ents); err != nil {
		return db.ReadExchangesResult{}, err
	}

	exchanges := make([]exchange.Exchange, len(ents))
	for i, e := range ents {
		exchanges[i] = e.ToModel()
	}

	return db.ReadExchangesResult{
		Exchanges: exchanges,
	}, nil
}

// UpdateExchanges will update the exchanges in the database.
func (a *Activities) UpdateExchanges(
	ctx context.Context,
	params db.UpdateExchangesParams,
) (db.UpdateExchangesResult, error) {
	ents := make([]interface{}, len(params.Exchanges))
	for i, e := range params.Exchanges {
		ents[i] = entities.ExchangeFromModel(e)
	}

	for _, e := range params.Exchanges {
		filter := bson.M{"name": e.Name}
		_, err := a.client.Collection(CollectionName).ReplaceOne(ctx, filter, entities.ExchangeFromModel(e))
		if err != nil {
			return db.UpdateExchangesResult{}, err
		}
	}

	return db.UpdateExchangesResult{}, nil
}

// DeleteExchanges will delete the exchanges from the database.
func (a *Activities) DeleteExchanges(
	ctx context.Context,
	params db.DeleteExchangesParams,
) (db.DeleteExchangesResult, error) {
	filter := bson.M{"name": bson.M{"$in": params.Names}}
	_, err := a.client.Collection(CollectionName).DeleteMany(ctx, filter)
	return db.DeleteExchangesResult{}, err
}
