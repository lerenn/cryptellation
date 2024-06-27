package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var exchangesCmd = &cobra.Command{
	Use:     "exchanges",
	Aliases: []string{"c"},
	Short:   "Manipulate exchanges service",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if err := executeParentPersistentPreRuns(cmd, args); err != nil {
			return err
		}

		return nil
	},
}

var exchangesReadCmd = &cobra.Command{
	Use:     "read",
	Aliases: []string{"r"},
	Short:   "Read exchanges from service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		list, err := globalClient.Exchanges().Read(context.Background(), "binance")
		if err != nil {
			return err
		}

		fmt.Println(list)
		return nil
	},
}

var exchangesInfoCmd = &cobra.Command{
	Use:     "service",
	Aliases: []string{"service"},
	Short:   "Read info from exchanges service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return displayServiceInfo(globalClient.Exchanges())
	},
}

func initExchanges(rootCmd *cobra.Command) {
	exchangesCmd.AddCommand(exchangesReadCmd)
	exchangesCmd.AddCommand(exchangesInfoCmd)
	rootCmd.AddCommand(exchangesCmd)
}
