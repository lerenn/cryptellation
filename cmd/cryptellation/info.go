package main

import (
	"github.com/lerenn/cryptellation/v1/clients/go/worker/client"
	"github.com/lerenn/cryptellation/v1/pkg/telemetry"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"i"},
	Short:   "Read info from worker",
	PersistentPreRunE: func(_ *cobra.Command, _ []string) (err error) {
		// Create cryptellation client
		cryptellationClient, err = client.NewClient()
		return err
	},
	PersistentPostRun: func(cmd *cobra.Command, _ []string) {
		cryptellationClient.Close(cmd.Context())
	},
	RunE: func(cmd *cobra.Command, _ []string) error {
		res, err := cryptellationClient.Info(cmd.Context())
		if err != nil {
			return err
		}

		telemetry.L(cmd.Context()).Infof("Version: %s", res.Version)

		return nil
	},
}
