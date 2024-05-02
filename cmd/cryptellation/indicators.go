package main

import (
	"fmt"

	indicators "github.com/lerenn/cryptellation/svc/indicators/clients/go/nats"
	"github.com/spf13/cobra"
)

var (
	indicatorsClient indicators.Client
)

var indicatorsCmd = &cobra.Command{
	Use:     "indicators",
	Aliases: []string{"c"},
	Short:   "Manipulate indicators service",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if err := executeParentPersistentPreRuns(cmd, args); err != nil {
			return err
		}

		indicatorsClient, err = indicators.NewClient(globalConfig)
		if err != nil {
			return fmt.Errorf("error when creating new indicators client: %w", err)
		}

		return nil
	},
}

var indicatorsInfoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"info"},
	Short:   "Read info from indicators service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return displayServiceInfo(indicatorsClient)
	},
}

func initIndicators(rootCmd *cobra.Command) {
	indicatorsCmd.AddCommand(indicatorsInfoCmd)
	rootCmd.AddCommand(indicatorsCmd)
}
