package main

import (
	"github.com/lerenn/cryptellation/v1/pkg/domains/ticks/activities/exchanges/live"
	"github.com/lerenn/cryptellation/v1/pkg/domains/ticks/workflows"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func registerTicksWorkflowsAndActivities(w worker.Worker, temporalClient client.Client) error {
	// Create exchange adapter
	exchanges, err := live.New(temporalClient)
	if err != nil {
		return err
	}
	exchanges.Register(w)

	// Create domain core
	domain := workflows.New(exchanges)
	domain.Register(w)

	return nil
}
