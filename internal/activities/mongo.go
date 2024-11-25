package activities

import (
	"context"
	"fmt"

	"github.com/lerenn/cryptellation/v1/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Mongo represents a Mongo client.
type Mongo struct {
	Client       *mongo.Client
	DatabaseName string
}

// NewMongo creates a new Mongo client.
func NewMongo(ctx context.Context, c config.Mongo) (Mongo, error) {
	if err := c.Validate(); err != nil {
		return Mongo{}, fmt.Errorf("loading mongo config: %w", err)
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(c.ConnectionString))
	if err != nil {
		return Mongo{}, fmt.Errorf("connecting to mongo: %w", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return Mongo{}, fmt.Errorf("pinging mongo: %w", err)
	}

	return Mongo{
		Client:       client,
		DatabaseName: c.Database,
	}, nil
}

// Collection returns the collection with the given name.
func (a *Mongo) Collection(collectionName string) *mongo.Collection {
	return a.Client.Database(a.DatabaseName).Collection(collectionName)
}

// DropDatabase drops the database.
func (a *Mongo) DropDatabase(ctx context.Context) error {
	return a.Client.Database(a.DatabaseName).Drop(ctx)
}
