package daemon

import (
	"context"

	"github.com/lerenn/cryptellation/internal/core/exchanges/ports/db"
	exchangesPort "github.com/lerenn/cryptellation/internal/core/exchanges/ports/exchanges"
	"github.com/lerenn/cryptellation/internal/infra/exchanges/db/sql"
	exchangesAdapter "github.com/lerenn/cryptellation/internal/infra/exchanges/exchanges"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/config/otel"
)

type adapters struct {
	db            db.Port
	exchanges     exchangesPort.Port
	otelExporters otel.Exporters
}

func newAdapters(ctx context.Context) (adapters, error) {
	// Init database client
	db, err := sql.New(config.LoadSQLConfigFromEnv())
	if err != nil {
		return adapters{}, err
	}

	// Init exchanges connections
	exchanges, err := exchangesAdapter.New(config.LoadExchangesConfigFromEnv())
	if err != nil {
		return adapters{}, err
	}

	// Init opentelemetry
	otelExporters, err := otel.NewExporters(ctx, "cryptellation-exchanges")
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
