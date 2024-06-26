package daemon

import (
	"context"

	asyncapipkg "github.com/lerenn/cryptellation/pkg/asyncapi"
	"github.com/lerenn/cryptellation/pkg/config"
	candlesticks "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"
	candlesticksNats "github.com/lerenn/cryptellation/svc/candlesticks/clients/go/nats"
	mongo "github.com/lerenn/cryptellation/svc/forwardtests/internal/adapters/db/mongo"
	"github.com/lerenn/cryptellation/svc/forwardtests/internal/app/ports/db"
)

type adapters struct {
	db           db.Port
	candlesticks candlesticks.Client
}

func newAdapters(ctx context.Context) (adapters, error) {
	// Init database client
	db, err := mongo.New(ctx, config.LoadMongo(nil))
	if err != nil {
		return adapters{}, err
	}

	// Init candlesticks client
	cdsClient, err := candlesticksNats.NewClient(
		config.LoadNATS(),
		candlesticksNats.WithLogger(asyncapipkg.LoggerWrapper{}),
		candlesticksNats.WithName("forwardtests"))
	if err != nil {
		return adapters{}, err
	}
	cachedCdsClient := candlesticks.NewCachedClient(cdsClient, candlesticks.DefaultCacheParameters())

	return adapters{
		db:           db,
		candlesticks: cachedCdsClient,
	}, nil
}

func (a adapters) Close(ctx context.Context) {
	a.candlesticks.Close(ctx)
}
