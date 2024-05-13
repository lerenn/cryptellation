package main

import (
	"github.com/spf13/cobra"
)

var indicatorsCmd = &cobra.Command{
	Use:     "indicators",
	Aliases: []string{"c"},
	Short:   "Manipulate indicators service",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if err := executeParentPersistentPreRuns(cmd, args); err != nil {
			return err
		}

		return nil
	},
}

var indicatorsInfoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"info"},
	Short:   "Read info from indicators service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return displayServiceInfo(globalClient.Indicators())
	},
}

func initIndicators(rootCmd *cobra.Command) {
	indicatorsCmd.AddCommand(indicatorsInfoCmd)
	rootCmd.AddCommand(indicatorsCmd)
}
