package main

import (
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
	RunE: func(cmd *cobra.Command, args []string) error {
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
			panic(err)
		}
		defer temporalClient.Close()

		// Create a worker
		w := worker.New(temporalClient, api.WorkerTaskQueueName, worker.Options{})

		// Register workflows
		w.RegisterWorkflowWithOptions(ServiceInfoWorkflow, workflow.RegisterOptions{
			Name: api.ServiceInfoWorkflowName,
		})

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

// ServiceInfoWorkflow returns the service information.
func ServiceInfoWorkflow(ctx workflow.Context) (api.ServiceInfoWorkflowResult, error) {
	return api.ServiceInfoWorkflowResult{
		Version: version.Version(),
	}, nil
}
