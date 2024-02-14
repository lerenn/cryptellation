package daemon

import (
	"context"

	asyncapipkg "github.com/lerenn/cryptellation/pkg/asyncapi"
	"github.com/lerenn/cryptellation/pkg/config"
	sql "github.com/lerenn/cryptellation/svc/backtests/internal/adapters/db/sql"
	natsBacktests "github.com/lerenn/cryptellation/svc/backtests/internal/adapters/events/nats"
	"github.com/lerenn/cryptellation/svc/backtests/internal/app/ports/db"
	"github.com/lerenn/cryptellation/svc/backtests/internal/app/ports/events"
	candlesticks "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"
	candlesticksNats "github.com/lerenn/cryptellation/svc/candlesticks/clients/go/nats"
)

type adapters struct {
	db           db.Port
	events       events.Port
	candlesticks candlesticks.Client
}

func newAdapters(ctx context.Context) (adapters, error) {
	// Init database client
	db, err := sql.New(config.LoadSQL())
	if err != nil {
		return adapters{}, err
	}

	// Init Events client
	events, err := natsBacktests.New(config.LoadNATS())
	if err != nil {
		return adapters{}, err
	}

	// Init candlesticks client
	candlesticks, err := candlesticksNats.NewClient(
		config.LoadNATS(), candlesticksNats.WithLogger(asyncapipkg.LoggerWrapper{}))
	if err != nil {
		return adapters{}, err
	}

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
