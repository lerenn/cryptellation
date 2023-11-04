package daemon

import (
	"context"

	"github.com/lerenn/cryptellation/internal/adapters/candlesticks/db/sql"
	"github.com/lerenn/cryptellation/internal/adapters/candlesticks/exchanges"
	"github.com/lerenn/cryptellation/internal/components/candlesticks/ports/db"
	exchangesIface "github.com/lerenn/cryptellation/internal/components/candlesticks/ports/exchanges"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/telemetry"
	"github.com/lerenn/cryptellation/pkg/telemetry/otel"
)

type adapters struct {
	db            db.Port
	exchanges     exchangesIface.Port
	otelExporters telemetry.Telemetry
}

func newAdapters(ctx context.Context) (adapters, error) {
	// Init database client
	db, err := sql.New(config.LoadSQLConfigFromEnv())
	if err != nil {
		return adapters{}, err
	}

	// Init exchanges connections
	exchanges, err := exchanges.New(config.LoadExchangesConfigFromEnv())
	if err != nil {
		return adapters{}, err
	}

	// Init opentelemetry
	otelExporters, err := otel.NewExporters(ctx, "cryptellation-candlesticks")
	if err != nil {
		return adapters{}, err
	}

	return adapters{
		db:            db,
		exchanges:     exchanges,
		otelExporters: otelExporters,
	}, nil
}

func (a adapters) Close(ctx context.Context) {
	a.otelExporters.Close(ctx)
}
