package daemon

import (
	"context"

	"github.com/lerenn/cryptellation/pkg/config"
	candlesticks "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"
	candlesticksNats "github.com/lerenn/cryptellation/svc/candlesticks/clients/go/nats"
	"github.com/lerenn/cryptellation/svc/indicators/internal/adapters/db/mongo"
	"github.com/lerenn/cryptellation/svc/indicators/internal/app/ports/db"
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
	candlesticks, err := candlesticksNats.NewClient(config.LoadNATS())
	if err != nil {
		return adapters{}, err
	}

	return adapters{
		db:           db,
		candlesticks: candlesticks,
	}, nil
}

func (a adapters) Close(ctx context.Context) {
	a.candlesticks.Close(ctx)
}
