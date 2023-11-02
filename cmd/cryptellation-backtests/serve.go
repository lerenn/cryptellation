package main

import (
	"context"

	"github.com/lerenn/cryptellation/cmd/cryptellation-backtests/daemon"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:     "serve",
	Aliases: []string{"s"},
	Short:   "Launch the service",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create a new daemon
		d, err := daemon.New(context.Background())
		if err != nil {
			return err
		}
		defer d.Close(context.Background())

		// Serve the daemon
		return d.Serve()
	},
}
