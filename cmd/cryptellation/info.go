package main

import (
	"github.com/lerenn/cryptellation/v1/pkg/telemetry"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"i"},
	Short:   "Read info from worker",
	RunE: func(cmd *cobra.Command, _ []string) error {
		res, err := cryptellationClient.Info(cmd.Context())
		if err != nil {
			return err
		}

		telemetry.L(cmd.Context()).Infof("Version: %s", res.Version)

		return nil
	},
}
