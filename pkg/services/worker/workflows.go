package worker

import (
	"context"

	api "github.com/lerenn/cryptellation/v1/api/worker/go"
	"github.com/lerenn/cryptellation/v1/pkg/config"
	"github.com/lerenn/cryptellation/v1/pkg/domains/backtests"
	backtestssql "github.com/lerenn/cryptellation/v1/pkg/domains/backtests/activities/db/sql"
	"github.com/lerenn/cryptellation/v1/pkg/domains/candlesticks"
	candlestickssql "github.com/lerenn/cryptellation/v1/pkg/domains/candlesticks/activities/db/sql"
	candlesticksexchagg "github.com/lerenn/cryptellation/v1/pkg/domains/candlesticks/activities/exchanges/aggregator"
	candlesticksbinance "github.com/lerenn/cryptellation/v1/pkg/domains/candlesticks/activities/exchanges/binance"
	"github.com/lerenn/cryptellation/v1/pkg/domains/exchanges"
	exchangessql "github.com/lerenn/cryptellation/v1/pkg/domains/exchanges/activities/db/sql"
	exchangesexchagg "github.com/lerenn/cryptellation/v1/pkg/domains/exchanges/activities/exchanges/aggregator"
	exchangesbinance "github.com/lerenn/cryptellation/v1/pkg/domains/exchanges/activities/exchanges/binance"
	"github.com/lerenn/cryptellation/v1/pkg/domains/forwardtests"
	forwardtestssql "github.com/lerenn/cryptellation/v1/pkg/domains/forwardtests/activities/db/sql"
	"github.com/lerenn/cryptellation/v1/pkg/domains/indicators"
	indicatorssql "github.com/lerenn/cryptellation/v1/pkg/domains/indicators/activities/db/sql"
	"github.com/lerenn/cryptellation/v1/pkg/domains/ticks"
	ticksexchagg "github.com/lerenn/cryptellation/v1/pkg/domains/ticks/activities/exchanges/aggregator"
	ticksbinance "github.com/lerenn/cryptellation/v1/pkg/domains/ticks/activities/exchanges/binance"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

// RegisterWorkflows registers all the workflows of the worker.
func RegisterWorkflows(ctx context.Context, w worker.Worker, temporalClient client.Client) error {
	// Register backtests workflows
	if err := registerBacktestsWorkflows(ctx, w); err != nil {
		return err
	}

	// Register candlesticks workflows
	if err := registerCandlesticksWorkflows(ctx, w); err != nil {
		return err
	}

	// Register exchanges workflows
	if err := registerExchangesWorkflows(ctx, w); err != nil {
		return err
	}

	// Register forwardtests workflows
	if err := registerForwardtestsWorkflows(ctx, w); err != nil {
		return err
	}

	// Register indicators workflows
	if err := registerIndicatorsWorkflows(ctx, w); err != nil {
		return err
	}

	// Register the service information workflow
	w.RegisterWorkflowWithOptions(ServiceInfo, workflow.RegisterOptions{
		Name: api.ServiceInfoWorkflowName,
	})

	// Register the ticks workflows
	return registerTicksWorkflows(w, temporalClient)
}

func registerBacktestsWorkflows(ctx context.Context, w worker.Worker) error {
	// Create database adapter
	db, err := backtestssql.New(ctx, config.LoadSQL(nil))
	if err != nil {
		return err
	}
	db.Register(w)

	// Create domain core
	domain := backtests.New(db)
	domain.Register(w)

	return nil
}

func registerCandlesticksWorkflows(ctx context.Context, w worker.Worker) error {
	// Create database adapter
	db, err := candlestickssql.New(ctx, config.LoadSQL(nil))
	if err != nil {
		return err
	}
	db.Register(w)

	// Create a binance exchange adapter
	binance, err := candlesticksbinance.New()
	if err != nil {
		return err
	}

	// Create exchange adapter aggregator
	exchanges := candlesticksexchagg.New(binance)
	exchanges.Register(w)

	// Create domain core
	domain := candlesticks.New(db, exchanges)
	domain.Register(w)

	return nil
}

func registerExchangesWorkflows(ctx context.Context, w worker.Worker) error {
	// Create database adapter
	db, err := exchangessql.New(ctx, config.LoadSQL(nil))
	if err != nil {
		return err
	}
	db.Register(w)

	// Create a binance exchange adapter
	binance, err := exchangesbinance.New()
	if err != nil {
		return err
	}

	// Create exchange adapter aggregator
	exchs := exchangesexchagg.New(binance)
	exchs.Register(w)

	// Create domain core
	domain := exchanges.New(db, exchs)
	domain.Register(w)

	return nil
}

func registerForwardtestsWorkflows(ctx context.Context, w worker.Worker) error {
	// Create database adapter
	db, err := forwardtestssql.New(ctx, config.LoadSQL(nil))
	if err != nil {
		return err
	}
	db.Register(w)

	// Create domain core
	domain := forwardtests.New(db)
	domain.Register(w)

	return nil
}

func registerIndicatorsWorkflows(ctx context.Context, w worker.Worker) error {
	// Create database adapter
	db, err := indicatorssql.New(ctx, config.LoadSQL(nil))
	if err != nil {
		return err
	}
	db.Register(w)

	// Create domain core
	domain := indicators.New(db)
	domain.Register(w)

	return nil
}

func registerTicksWorkflows(w worker.Worker, temporalClient client.Client) error {
	// Create a binance exchange adapter
	binance, err := ticksbinance.New(temporalClient)
	if err != nil {
		return err
	}

	// Create exchange adapter aggregator
	exchanges := ticksexchagg.New(binance)
	exchanges.Register(w)

	// Create domain core
	domain := ticks.New(exchanges)
	domain.Register(w)

	return nil
}
