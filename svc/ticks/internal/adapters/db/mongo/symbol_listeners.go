package mongo

import (
	"context"

	"github.com/lerenn/cryptellation/svc/ticks/internal/adapters/db/mongo/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// CollectionName is the name of the collection in the database
	CollectionName = "symbol_listeners"
)

func (a *Adapter) IncrementSymbolListenerSubscribers(ctx context.Context, exchange, pair string) (int64, error) {
	var result entities.SymbolListener

	err := a.client.Collection(CollectionName).FindOneAndUpdate(
		ctx,
		bson.M{
			"exchange": exchange,
			"pair":     pair,
		},
		bson.M{
			"$inc": bson.M{
				"subscribers": 1,
			},
		},
		options.FindOneAndUpdate().SetUpsert(true),
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&result)

	return result.Subscribers, err
}

func (a *Adapter) DecrementSymbolListenerSubscribers(ctx context.Context, exchange, pair string) (int64, error) {
	var result entities.SymbolListener

	err := a.client.Collection(CollectionName).FindOneAndUpdate(
		ctx,
		bson.M{
			"exchange": exchange,
			"pair":     pair,
		},
		bson.M{
			"$inc": bson.M{
				"subscribers": -1,
			},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&result)

	return result.Subscribers, err
}

func (a *Adapter) GetSymbolListenerSubscribers(ctx context.Context, exchange, pair string) (int64, error) {
	var result entities.SymbolListener

	err := a.client.Collection(CollectionName).FindOne(
		ctx,
		bson.M{
			"exchange": exchange,
			"pair":     pair,
		},
	).Decode(&result)

	return result.Subscribers, err
}

func (a *Adapter) ClearSymbolListenerSubscribers(ctx context.Context, exchange, pair string) error {
	_, err := a.client.Collection(CollectionName).DeleteOne(
		ctx,
		bson.M{
			"exchange": exchange,
			"pair":     pair,
		},
	)
	return err
}

func (a *Adapter) ClearAllSymbolListenersCount(ctx context.Context) error {
	_, err := a.client.Collection(CollectionName).UpdateMany(
		ctx,
		bson.M{},
		bson.M{
			"$set": bson.M{
				"subscribers": 0,
			},
		},
	)
	return err
}
