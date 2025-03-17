package main

import (
	"fmt"
	"time"

	api "github.com/lerenn/cryptellation/v1/api/worker/go"
	"github.com/lerenn/cryptellation/v1/pkg/config"
	"github.com/lerenn/cryptellation/v1/pkg/health"
	"github.com/lerenn/cryptellation/v1/pkg/services/worker"
	"github.com/lerenn/cryptellation/v1/pkg/telemetry"
	"github.com/lerenn/cryptellation/v1/pkg/temporal/activities"
	"github.com/spf13/cobra"
	"go.temporal.io/sdk/client"
	temporalwk "go.temporal.io/sdk/worker"
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
		w := temporalwk.New(temporalClient, api.WorkerTaskQueueName, temporalwk.Options{})

		// Register common activities
		w.RegisterActivity(activities.NewActivities(temporalClient))

		// Register workflows
		if err := worker.RegisterWorkflows(cmd.Context(), w, temporalClient); err != nil {
			return err
		}

		// Mark as ready
		// TODO(#55): Improve this with a better way to mark as ready
		go func() {
			time.Sleep(time.Second * 3)
			h.Ready(true)
		}()
		defer h.Ready(false)

		// Run worker
		return w.Run(temporalwk.InterruptCh())
	},
}
