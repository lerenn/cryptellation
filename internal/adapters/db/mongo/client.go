package mongo

import (
	"context"
	"fmt"

	"cryptellation/pkg/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Client struct {
	Client       *mongo.Client
	DatabaseName string
}

func NewClient(ctx context.Context, c config.Mongo) (Client, error) {
	if err := c.Validate(); err != nil {
		return Client{}, fmt.Errorf("loading mongo config: %w", err)
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(c.ConnectionString))
	if err != nil {
		return Client{}, fmt.Errorf("connecting to mongo: %w", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return Client{}, fmt.Errorf("pinging mongo: %w", err)
	}

	return Client{
		Client:       client,
		DatabaseName: c.Database,
	}, nil
}

func (a *Client) Collection(collectionName string) *mongo.Collection {
	return a.Client.Database(a.DatabaseName).Collection(collectionName)
}

func (a *Client) DropDatabase(ctx context.Context) error {
	return a.Client.Database(a.DatabaseName).Drop(ctx)
}
