package mongo

import (
	"context"

	"github.com/lerenn/cryptellation/exchanges/internal/adapters/db/mongo/entities"
	"github.com/lerenn/cryptellation/exchanges/pkg/exchange"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	// CollectionName is the name of the collection in the database
	CollectionName = "exchanges"
)

func (a *Adapter) CreateExchanges(ctx context.Context, exchanges ...exchange.Exchange) error {
	ents := make([]interface{}, len(exchanges))
	for i, e := range exchanges {
		ents[i] = entities.ExchangeFromModel(e)
	}

	_, err := a.client.Collection(CollectionName).InsertMany(ctx, ents)
	return err
}

func (a *Adapter) ReadExchanges(ctx context.Context, names ...string) ([]exchange.Exchange, error) {
	filter := bson.M{}
	if len(names) > 0 {
		filter["name"] = bson.M{"$in": names}
	}

	cur, err := a.client.Collection(CollectionName).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var ents []entities.Exchange
	if err := cur.All(ctx, &ents); err != nil {
		return nil, err
	}

	exchanges := make([]exchange.Exchange, len(ents))
	for i, e := range ents {
		exchanges[i] = e.ToModel()
	}

	return exchanges, nil
}

func (a *Adapter) UpdateExchanges(ctx context.Context, exchanges ...exchange.Exchange) error {
	ents := make([]interface{}, len(exchanges))
	for i, e := range exchanges {
		ents[i] = entities.ExchangeFromModel(e)
	}

	for _, e := range exchanges {
		filter := bson.M{"name": e.Name}
		_, err := a.client.Collection(CollectionName).ReplaceOne(ctx, filter, entities.ExchangeFromModel(e))
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *Adapter) DeleteExchanges(ctx context.Context, names ...string) error {
	filter := bson.M{"name": bson.M{"$in": names}}
	_, err := a.client.Collection(CollectionName).DeleteMany(ctx, filter)
	return err
}
