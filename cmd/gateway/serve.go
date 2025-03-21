package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lerenn/cryptellation/v1/clients/go/worker/client"
	"github.com/lerenn/cryptellation/v1/pkg/health"
	"github.com/lerenn/cryptellation/v1/pkg/services/gateway"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:     "serve",
	Aliases: []string{"s"},
	Short:   "Launch the server",
	RunE: func(cmd *cobra.Command, _ []string) error {
		// Init and serve health server
		// NOTE: health OK, but not-ready yet
		h, err := health.NewHealth(cmd.Context())
		if err != nil {
			return err
		}
		go h.HTTPServe(cmd.Context())

		// Create cryptellation client
		cryptellationClient, err := client.NewClient()
		if err != nil {
			return err
		}
		defer cryptellationClient.Close(cmd.Context())

		// Create server
		server := gateway.NewServer(cryptellationClient)

		// Create router and set routes
		router := gin.Default()
		gateway.RegisterHandlers(router.Group("v1"), server)

		// Mark as ready
		// TODO(#54): Improve this with a better way to mark as ready
		go func() {
			time.Sleep(time.Second * 3)
			h.Ready(true)
		}()
		defer h.Ready(false)

		// Run worker
		return router.Run(":8080")
	},
}
