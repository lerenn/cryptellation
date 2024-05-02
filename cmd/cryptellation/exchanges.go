package main

import (
	"context"
	"fmt"

	exchanges "github.com/lerenn/cryptellation/svc/exchanges/clients/go/nats"
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

		exchangesClient, err = exchanges.NewClient(globalConfig)
		if err != nil {
			return fmt.Errorf("error when creating new candlesticks client: %w", err)
		}

		return nil
	},
}

var exchangesReadCmd = &cobra.Command{
	Use:     "read",
	Aliases: []string{"r"},
	Short:   "Read exchanges from service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		list, err := exchangesClient.Read(context.Background(), "binance")
		if err != nil {
			return err
		}

		fmt.Println(list)
		return nil
	},
}

var exchangesInfoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"info"},
	Short:   "Read info from exchanges service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return displayServiceInfo(exchangesClient)
	},
}

func initExchanges(rootCmd *cobra.Command) {
	exchangesCmd.AddCommand(exchangesReadCmd)
	exchangesCmd.AddCommand(exchangesInfoCmd)
	rootCmd.AddCommand(exchangesCmd)
}
