package main

import (
	"context"

	"github.com/lerenn/cryptellation/v1/internal/domains/exchanges/activities/db/mongo"
	"github.com/lerenn/cryptellation/v1/internal/domains/exchanges/activities/exchanges/live"
	"github.com/lerenn/cryptellation/v1/internal/domains/exchanges/workflows"
	"github.com/lerenn/cryptellation/v1/pkg/config"
	"go.temporal.io/sdk/worker"
)

func registerExchangesWorkflowsAndActivities(ctx context.Context, w worker.Worker) error {
	// Create database activities
	dbActivities, err := mongo.New(ctx, config.LoadMongo(nil))
	if err != nil {
		return err
	}
	dbActivities.Register(w)

	// Create exchange activities
	exchangesActivities, err := live.New()
	if err != nil {
		return err
	}
	exchangesActivities.Register(w)

	// Create exchanges domain
	exchanges := workflows.New(dbActivities, exchangesActivities)

	// Register workflows
	exchanges.Register(w)

	return nil
}
