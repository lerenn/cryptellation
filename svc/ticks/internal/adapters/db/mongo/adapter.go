package mongo

import (
	"context"

	mongoutil "github.com/lerenn/cryptellation/pkg/adapters/db/mongo"
	"github.com/lerenn/cryptellation/pkg/config"
	port "github.com/lerenn/cryptellation/svc/ticks/internal/app/ports/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ port.Port = (*Adapter)(nil)

type Adapter struct {
	client mongoutil.Client
}

func New(ctx context.Context, c config.Mongo) (*Adapter, error) {
	// Create embedded database access
	db, err := mongoutil.NewClient(ctx, c)
	if err != nil {
		return nil, err
	}

	// Create a structure
	a := &Adapter{
		client: db,
	}

	// Create indexes
	return a, a.CreateIndexes()
}

func (a *Adapter) CreateIndexes() error {
	_, err := a.client.
		Collection(CollectionName).
		Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "exchange", Value: 1},
			{Key: "pair", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})

	return err
}

func (a *Adapter) Reset(ctx context.Context) error {
	if err := a.client.DropDatabase(ctx); err != nil {
		return err
	}

	return a.CreateIndexes()
}
