package main

import (
	"context"
	"fmt"

	exchanges "github.com/digital-feather/cryptellation/internal/exchanges/ctrl/nats"
	"github.com/spf13/cobra"
)

var (
	exchangesClient exchanges.Client
)

var exchangesCmd = &cobra.Command{
	Use:     "exchanges",
	Aliases: []string{"c"},
	Short:   "Manipulate exchanges service",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if err := executeParentPersistentPreRuns(cmd, args); err != nil {
			return err
		}

		exchangesClient, err = exchanges.New(globalNATSConfig)
		return err
	},
}

var exchangesReadCmd = &cobra.Command{
	Use:     "read",
	Aliases: []string{"r"},
	Short:   "Read exchanges from service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		list, err := exchangesClient.ReadExchanges(context.Background(), "binance")
		if err != nil {
			return err
		}

		fmt.Println(list)
		return nil
	},
}

func initExchanges(rootCmd *cobra.Command) {
	exchangesCmd.AddCommand(exchangesReadCmd)
	rootCmd.AddCommand(exchangesCmd)
}
