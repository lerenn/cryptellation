package mongo

import (
	"context"
	"fmt"

	"github.com/lerenn/cryptellation/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Client struct {
	client       *mongo.Client
	databaseName string
}

func NewClient(c config.Mongo) (Client, error) {
	if err := c.Validate(); err != nil {
		return Client{}, fmt.Errorf("loading mongo config: %w", err)
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(c.ConnectionString))
	if err != nil {
		return Client{}, fmt.Errorf("connecting to mongo: %w", err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return Client{}, fmt.Errorf("pinging mongo: %w", err)
	}

	return Client{
		client:       client,
		databaseName: c.Database,
	}, nil
}

func (a *Client) Collection(collectionName string) *mongo.Collection {
	return a.client.Database(a.databaseName).Collection(collectionName)
}

func (a *Client) DropDatabase(ctx context.Context) error {
	return a.client.Database(a.databaseName).Drop(ctx)
}
