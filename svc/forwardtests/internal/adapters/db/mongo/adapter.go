package mongo

import (
	"context"

	mongoutil "github.com/lerenn/cryptellation/internal/adapters/db/mongo"

	"github.com/lerenn/cryptellation/pkg/config"

	port "github.com/lerenn/cryptellation/svc/forwardtests/internal/app/ports/db"

	"go.mongodb.org/mongo-driver/mongo"
)

var _ port.Port = (*Adapter)(nil)

type Adapter struct {
	client mongoutil.Client
}

func New(ctx context.Context, c config.Mongo) (*Adapter, error) {
	// Create embedded database access
	db, err := mongoutil.NewClient(ctx, c)

	// Return database access
	return &Adapter{
		client: db,
	}, err
}

func (a *Adapter) CreateIndexes() error {
	index := mongo.IndexModel{Keys: map[string]int{"updated_at": 1}}

	// Create indexes
	_, err := a.client.
		Collection(CollectionName).
		Indexes().
		CreateOne(context.Background(), index)
	if err != nil {
		return err
	}

	return nil
}

func (a *Adapter) Reset(ctx context.Context) error {
	if err := a.client.DropDatabase(ctx); err != nil {
		return err
	}

	return a.CreateIndexes()
}
