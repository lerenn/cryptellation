package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lerenn/cryptellation/v1/pkg/health"
	"github.com/lerenn/cryptellation/v1/web/ui"
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

		router := gin.Default()
		ui.AddRoutes(router)

		// Mark as ready
		// TODO: improve this
		go func() {
			time.Sleep(time.Second * 3)
			h.Ready(true)
		}()
		defer h.Ready(false)

		// Run worker
		return router.Run(":8080")
	},
}
