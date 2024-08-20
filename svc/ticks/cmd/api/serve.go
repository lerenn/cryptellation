package main

import (
	"github.com/lerenn/cryptellation/ticks/cmd/api/daemon"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:     "serve",
	Aliases: []string{"s"},
	Short:   "Launch the service",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create a new daemon
		d, err := daemon.New(cmd.Context())
		if err != nil {
			return err
		}
		defer d.Close(cmd.Context())

		// Serve the daemon
		return d.Serve(cmd.Context())
	},
}
