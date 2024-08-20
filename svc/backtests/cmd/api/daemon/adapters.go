package daemon

import (
	"context"

	asyncapipkg "github.com/lerenn/cryptellation/internal/asyncapi"

	"github.com/lerenn/cryptellation/pkg/config"

	mongo "github.com/lerenn/cryptellation/svc/backtests/internal/adapters/db/mongo"
	natsBacktests "github.com/lerenn/cryptellation/svc/backtests/internal/adapters/events/nats"
	"github.com/lerenn/cryptellation/svc/backtests/internal/app/ports/db"
	"github.com/lerenn/cryptellation/svc/backtests/internal/app/ports/events"

	candlesticks "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"
	candlestickscache "github.com/lerenn/cryptellation/svc/candlesticks/clients/go/cache"
	candlesticksnats "github.com/lerenn/cryptellation/svc/candlesticks/clients/go/nats"
	candlesticksretry "github.com/lerenn/cryptellation/svc/candlesticks/clients/go/retry"
)

type adapters struct {
	db           db.Port
	events       events.Port
	candlesticks candlesticks.Client
}

func newAdapters(ctx context.Context) (adapters, error) {
	// Init database client
	db, err := mongo.New(ctx, config.LoadMongo(nil))
	if err != nil {
		return adapters{}, err
	}

	// Init Events client
	events, err := natsBacktests.New(config.LoadNATS())
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
		events:       events,
		candlesticks: candlesticks,
	}, nil
}

func (a adapters) Close(ctx context.Context) {
	a.candlesticks.Close(ctx)
	a.events.Close(ctx)
}
