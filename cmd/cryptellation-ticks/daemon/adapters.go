package daemon

import (
	"context"

	"github.com/lerenn/cryptellation/internal/adapters/ticks/db/sql"
	natsAdapter "github.com/lerenn/cryptellation/internal/adapters/ticks/events/nats"
	exchangesAdapter "github.com/lerenn/cryptellation/internal/adapters/ticks/exchanges"
	"github.com/lerenn/cryptellation/internal/components/ticks/ports/db"
	"github.com/lerenn/cryptellation/internal/components/ticks/ports/events"
	exchangesPort "github.com/lerenn/cryptellation/internal/components/ticks/ports/exchanges"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/telemetry"
	"github.com/lerenn/cryptellation/pkg/telemetry/otel"
)

type adapters struct {
	db            db.Port
	events        events.Port
	exchanges     exchangesPort.Port
	otelExporters telemetry.Telemetry
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

	// Init Events client
	events, err := natsAdapter.New(config.LoadNATSConfigFromEnv())
	if err != nil {
		return adapters{}, err
	}

	// Init opentelemetry
	otelExporters, err := otel.NewExporters(ctx, "cryptellation-ticks")
	if err != nil {
		return adapters{}, err
	}

	return adapters{
		db:            db,
		events:        events,
		exchanges:     exchanges,
		otelExporters: otelExporters,
	}, nil
}

func (a adapters) Close(ctx context.Context) {
	a.events.Close(ctx)
	a.otelExporters.Close(ctx)
}
