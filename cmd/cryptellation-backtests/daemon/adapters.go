package daemon

import (
	"context"

	client "github.com/lerenn/cryptellation/clients/go"
	natsClient "github.com/lerenn/cryptellation/clients/go/nats"
	"github.com/lerenn/cryptellation/internal/adapters/backtests/db/sql"
	natsAdapter "github.com/lerenn/cryptellation/internal/adapters/backtests/events/nats"
	"github.com/lerenn/cryptellation/internal/components/backtests/ports/db"
	"github.com/lerenn/cryptellation/internal/components/backtests/ports/events"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/telemetry"
	"github.com/lerenn/cryptellation/pkg/telemetry/otel"
)

type adapters struct {
	db            db.Port
	events        events.Port
	candlesticks  client.Candlesticks
	otelExporters telemetry.Telemetry
}

func newAdapters(ctx context.Context) (adapters, error) {
	// Init database client
	db, err := sql.New(config.LoadSQLConfigFromEnv())
	if err != nil {
		return adapters{}, err
	}

	// Init Events client
	events, err := natsAdapter.New(config.LoadNATSConfigFromEnv())
	if err != nil {
		return adapters{}, err
	}

	// Init candlesticks client
	candlesticks, err := natsClient.NewCandlesticks(config.LoadNATSConfigFromEnv())
	if err != nil {
		return adapters{}, err
	}

	// Init opentelemetry
	otelExporters, err := otel.NewExporters(ctx, "cryptellation-backtests")
	if err != nil {
		return adapters{}, err
	}

	return adapters{
		db:            db,
		events:        events,
		candlesticks:  candlesticks,
		otelExporters: otelExporters,
	}, nil
}

func (a adapters) Close(ctx context.Context) {
	a.otelExporters.Close(ctx)
	a.candlesticks.Close(ctx)
	a.events.Close(ctx)
}
