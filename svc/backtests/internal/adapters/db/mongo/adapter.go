package backtests

import (
	"context"

	mongoutil "cryptellation/internal/adapters/db/mongo"
	"cryptellation/pkg/config"

	port "cryptellation/svc/backtests/internal/app/ports/db"
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
