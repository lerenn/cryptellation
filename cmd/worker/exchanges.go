package main

import (
	"context"

	"github.com/lerenn/cryptellation/v1/pkg/config"
	"github.com/lerenn/cryptellation/v1/pkg/domains/exchanges/activities/db/mongo"
	"github.com/lerenn/cryptellation/v1/pkg/domains/exchanges/activities/exchanges/live"
	"github.com/lerenn/cryptellation/v1/pkg/domains/exchanges/workflows"
	"go.temporal.io/sdk/worker"
)

func registerExchangesWorkflowsAndActivities(ctx context.Context, w worker.Worker) error {
	// Create database adapter
	db, err := mongo.New(ctx, config.LoadMongo(nil))
	if err != nil {
		return err
	}
	db.Register(w)

	// Create exchange adapter
	exchanges, err := live.New()
	if err != nil {
		return err
	}
	exchanges.Register(w)

	// Create domain core
	domain := workflows.New(db, exchanges)
	domain.Register(w)

	return nil
}
