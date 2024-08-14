package daemon

import (
	"context"

	"cryptellation/internal/config"

	mongo "cryptellation/svc/exchanges/internal/adapters/db/mongo"
	"cryptellation/svc/exchanges/internal/adapters/exchanges"
	"cryptellation/svc/exchanges/internal/app/ports/db"
	exchangesPort "cryptellation/svc/exchanges/internal/app/ports/exchanges"
)

type adapters struct {
	db        db.Port
	exchanges exchangesPort.Port
}

func newAdapters(ctx context.Context) (adapters, error) {
	// Init database client
	db, err := mongo.New(ctx, config.LoadMongo(nil))
	if err != nil {
		return adapters{}, err
	}

	// Init exchanges connections
	exchanges, err := exchanges.New()
	if err != nil {
		return adapters{}, err
	}

	return adapters{
		db:        db,
		exchanges: exchanges,
	}, nil
}

func (a adapters) Close(ctx context.Context) {
}
