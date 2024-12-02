package main

import (
	"fmt"

	"github.com/lerenn/cryptellation/v1/api"
	"github.com/spf13/cobra"
)

var (
	ticksListenExchangeFlag string
	ticksListenPairFlag     string
)

var ticksCmd = &cobra.Command{
	Use:     "ticks",
	Aliases: []string{"t"},
	Short:   "Manage ticks",
}

var ticksListenCmd = &cobra.Command{
	Use:     "listen",
	Aliases: []string{"l"},
	Short:   "Listen to ticks",
	RunE: func(cmd *cobra.Command, _ []string) error {
		res, err := cryptellationClient.ListenToTicks(cmd.Context(), api.RegisterForTicksListeningWorkflowParams{
			Exchange: ticksListenExchangeFlag,
			Pair:     ticksListenPairFlag,
			// CallbackWorkflow: TODO,
		})
		if err != nil {
			return err
		}

		fmt.Println(res)

		return nil
	},
}

func addTicksCommands() {
	ticksListenCmd.Flags().StringVarP(&ticksListenExchangeFlag, "exchange", "e", "binance", "Exchange")
	ticksListenCmd.Flags().StringVarP(&ticksListenPairFlag, "pair", "p", "BTC-USDT", "Pair")
	ticksCmd.AddCommand(ticksListenCmd)
	rootCmd.AddCommand(ticksCmd)
}
