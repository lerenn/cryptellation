package main

import (
	"fmt"

	backtests "github.com/lerenn/cryptellation/svc/backtests/clients/go/nats"
	"github.com/spf13/cobra"
)

var (
	backtestsClient backtests.Client
)

var backtestsCmd = &cobra.Command{
	Use:     "backtests",
	Aliases: []string{"c"},
	Short:   "Manipulate backtests service",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if err := executeParentPersistentPreRuns(cmd, args); err != nil {
			return err
		}

		backtestsClient, err = backtests.NewClient(globalConfig)
		if err != nil {
			return fmt.Errorf("error when creating new backtests client: %w", err)
		}

		return nil
	},
}

var backtestsInfoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"info"},
	Short:   "Read info from backtests service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return displayServiceInfo(backtestsClient)
	},
}

func initBacktests(rootCmd *cobra.Command) {
	backtestsCmd.AddCommand(backtestsInfoCmd)
	rootCmd.AddCommand(backtestsCmd)
}
