package daemon

import (
	"context"

	asyncapipkg "cryptellation/internal/asyncapi"
	"cryptellation/pkg/config"

	candlesticks "cryptellation/svc/candlesticks/clients/go"
	candlestickscache "cryptellation/svc/candlesticks/clients/go/cache"
	candlesticksnats "cryptellation/svc/candlesticks/clients/go/nats"
	candlesticksretry "cryptellation/svc/candlesticks/clients/go/retry"

	"cryptellation/svc/indicators/internal/adapters/db/mongo"
	"cryptellation/svc/indicators/internal/app/ports/db"
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
	candlesticks, err := candlesticksnats.New(
		config.LoadNATS(),
		candlesticksnats.WithLogger(asyncapipkg.LoggerWrapper{}),
		candlesticksnats.WithName("backtests"))
	if err != nil {
		return adapters{}, err
	}
	candlesticks = candlestickscache.New(candlesticks)
	candlesticks = candlesticksretry.New(candlesticks)

	return adapters{
		db:           db,
		candlesticks: candlesticks,
	}, nil
}

func (a adapters) Close(ctx context.Context) {
	a.candlesticks.Close(ctx)
}
