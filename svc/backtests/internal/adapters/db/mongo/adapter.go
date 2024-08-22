package backtests

import (
	"context"

	mongoutil "github.com/lerenn/cryptellation/internal/adapters/db/mongo"

	"github.com/lerenn/cryptellation/pkg/config"

	port "github.com/lerenn/cryptellation/svc/backtests/internal/app/ports/db"
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
	// No index
	return nil
}

func (a *Adapter) Reset(ctx context.Context) error {
	if err := a.client.DropDatabase(ctx); err != nil {
		return err
	}

	return a.CreateIndexes()
}
