package main

import (
	"fmt"
	"time"

	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/pkg/config"
	"github.com/lerenn/cryptellation/v1/pkg/health"
	"github.com/lerenn/cryptellation/v1/pkg/telemetry"
	"github.com/lerenn/cryptellation/v1/pkg/temporal/activities"
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
		var temporalClient client.Client
		for {
			temporalClient, err = client.Dial(client.Options{
				HostPort: temporalConfig.Address,
			})
			if err != nil {
				msg := fmt.Sprintf("cannot connect to temporal: %s", err)
				telemetry.L(cmd.Context()).Warning(msg)
				time.Sleep(3 * time.Second)
			} else {
				break
			}
		}
		defer temporalClient.Close()

		// Create a worker
		w := worker.New(temporalClient, api.WorkerTaskQueueName, worker.Options{})

		// Register common activities
		w.RegisterActivity(activities.NewActivities(temporalClient))

		// Register workflows
		if err := registerWorkflows(cmd.Context(), w, temporalClient); err != nil {
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

// ServiceInfo returns the service information.
func ServiceInfo(_ workflow.Context, _ api.ServiceInfoParams) (api.ServiceInfoResults, error) {
	return api.ServiceInfoResults{
		Version: version.Version(),
	}, nil
}
