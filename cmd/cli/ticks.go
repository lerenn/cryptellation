package main

import (
	"fmt"

	"github.com/lerenn/cryptellation/v1/api"
	"github.com/spf13/cobra"
	"go.temporal.io/sdk/workflow"
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
		return cryptellationClient.ListenToTicks(cmd.Context(),
			ticksListenExchangeFlag,
			ticksListenPairFlag,
			func(_ workflow.Context, params api.ListenToTicksCallbackWorkflowParams) error {
				fmt.Println(params.Tick.String())
				return nil
			},
		)
	},
}

func addTicksCommands() {
	ticksListenCmd.Flags().StringVarP(&ticksListenExchangeFlag, "exchange", "e", "binance", "Exchange")
	ticksListenCmd.Flags().StringVarP(&ticksListenPairFlag, "pair", "p", "BTC-USDT", "Pair")
	ticksCmd.AddCommand(ticksListenCmd)
	rootCmd.AddCommand(ticksCmd)
}
