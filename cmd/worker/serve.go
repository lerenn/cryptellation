package main

import (
	"context"
	"time"

	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/pkg/config"
	"github.com/lerenn/cryptellation/v1/pkg/health"
	"github.com/lerenn/cryptellation/v1/pkg/version"
	"github.com/spf13/cobra"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

var serveCmd = &cobra.Command{
	Use:     "serve",
	Aliases: []string{"s"},
	Short:   "Launch the service",
	RunE: func(cmd *cobra.Command, _ []string) error {
		// Init and serve health server
		// NOTE: health OK, but not-ready yet
		h, err := health.NewHealth(cmd.Context())
		if err != nil {
			return err
		}
		go h.HTTPServe(cmd.Context())

		// Load temporal configuration
		temporalConfig := config.LoadTemporal(nil)
		if err := temporalConfig.Validate(); err != nil {
			return err
		}

		// Load temporal client
		temporalClient, err := client.Dial(client.Options{
			HostPort: temporalConfig.Address,
		})
		if err != nil {
			return err
		}
		defer temporalClient.Close()

		// Create a worker
		w := worker.New(temporalClient, api.WorkerTaskQueueName, worker.Options{})

		// Register workflows
		if err := registerWorflowsAndActivities(cmd.Context(), w); err != nil {
			return err
		}

		// Mark as ready
		// TODO: improve this
		go func() {
			time.Sleep(time.Second * 3)
			h.Ready(true)
		}()
		defer h.Ready(false)

		// Run worker
		return w.Run(worker.InterruptCh())
	},
}

func registerWorflowsAndActivities(ctx context.Context, w worker.Worker) error {
	// Register candlesticks workflows
	if err := registerCandlesticksWorkflowsAndActivities(ctx, w); err != nil {
		return err
	}

	// Register exchanges workflows
	if err := registerExchangesWorkflowsAndActivities(ctx, w); err != nil {
		return err
	}

	// Register the service information workflow
	w.RegisterWorkflowWithOptions(ServiceInfo, workflow.RegisterOptions{
		Name: api.ServiceInfoWorkflowName,
	})

	return nil
}

// ServiceInfo returns the service information.
func ServiceInfo(_ workflow.Context, _ api.ServiceInfoParams) (api.ServiceInfoResult, error) {
	return api.ServiceInfoResult{
		Version: version.Version(),
	}, nil
}
