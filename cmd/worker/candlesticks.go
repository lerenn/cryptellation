package main

import (
	"context"

	"github.com/lerenn/cryptellation/v1/internal/domains/candlesticks/activities/db/mongo"
	"github.com/lerenn/cryptellation/v1/internal/domains/candlesticks/activities/exchanges/live"
	"github.com/lerenn/cryptellation/v1/internal/domains/candlesticks/workflows"
	"github.com/lerenn/cryptellation/v1/pkg/config"
	"go.temporal.io/sdk/worker"
)

func registerCandlesticksWorkflowsAndActivities(ctx context.Context, w worker.Worker) error {
	// Create database activity
	dbAdapter, err := mongo.New(ctx, config.LoadMongo(nil))
	if err != nil {
		return err
	}
	dbAdapter.Register(w)

	// Create exchange activity
	exchangesAdapter, err := live.New()
	if err != nil {
		return err
	}
	exchangesAdapter.Register(w)

	// Create candlesticks domain
	candlesticks := workflows.New(dbAdapter, exchangesAdapter)

	// Register workflows
	candlesticks.Register(w)

	return nil
}
