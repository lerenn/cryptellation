package main

import (
	"context"
	"fmt"

	client "github.com/lerenn/cryptellation/clients/go"
	"github.com/lerenn/cryptellation/clients/go/nats"
	"github.com/spf13/cobra"
)

var (
	exchangesClient client.Exchanges
)

var exchangesCmd = &cobra.Command{
	Use:     "exchanges",
	Aliases: []string{"c"},
	Short:   "Manipulate exchanges service",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if err := executeParentPersistentPreRuns(cmd, args); err != nil {
			return err
		}

		exchangesClient, err = nats.NewExchanges(globalNATSConfig)
		return err
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

func initExchanges(rootCmd *cobra.Command) {
	exchangesCmd.AddCommand(exchangesReadCmd)
	rootCmd.AddCommand(exchangesCmd)
}
