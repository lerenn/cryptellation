package main

import (
	"fmt"
	"time"

	api "github.com/lerenn/cryptellation/v1/api/worker/go"
	"github.com/lerenn/cryptellation/v1/clients/go/worker/client"
	"github.com/lerenn/cryptellation/v1/pkg/models/period"
	"github.com/spf13/cobra"
)

var (
	candlesticksListExchangeFlag string
	candlesticksListPairFlag     string
	candlesticksListPeriodFlag   string
	candlesticksListStartFlag    string
	candlesticksListEndFlag      string
)

var candlesticksCmd = &cobra.Command{
	Use:     "candlesticks",
	Aliases: []string{"c"},
	Short:   "List candlesticks",
	PersistentPreRunE: func(_ *cobra.Command, _ []string) (err error) {
		// Create cryptellation client
		cryptellationClient, err = client.NewClient()
		return err
	},
	PersistentPostRun: func(cmd *cobra.Command, _ []string) {
		cryptellationClient.Close(cmd.Context())
	},
}

var candlesticksListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List candlesticks",
	RunE: func(cmd *cobra.Command, _ []string) error {
		// Parse start date
		start, err := time.Parse(time.RFC3339, candlesticksListStartFlag)
		if err != nil {
			return err
		}

		// Parse end date
		end, err := time.Parse(time.RFC3339, candlesticksListEndFlag)
		if err != nil {
			return err
		}

		// Parse period
		per, err := period.FromString(candlesticksListPeriodFlag)
		if err != nil {
			return err
		}

		// Execute call
		res, err := cryptellationClient.ListCandlesticks(cmd.Context(), api.ListCandlesticksWorkflowParams{
			Exchange: candlesticksListExchangeFlag,
			Pair:     candlesticksListPairFlag,
			Period:   per,
			Start:    &start,
			End:      &end,
		})
		if err != nil {
			return err
		}

		switch {
		case jsonOutput:
			return displayJSON(cmd.Context(), res.List.ToArray())
		default:
			fmt.Println(res.List.String())
		}

		return nil
	},
}

func addCandlesticksCommands() {
	candlesticksListCmd.Flags().StringVarP(&candlesticksListExchangeFlag, "exchange", "e", "binance", "Exchange")
	candlesticksListCmd.Flags().StringVarP(&candlesticksListPairFlag, "pair", "p", "ETH-USDT", "Pair")
	candlesticksListCmd.Flags().StringVarP(&candlesticksListPeriodFlag, "period", "P", "H1", "Period")
	candlesticksListCmd.Flags().StringVarP(
		&candlesticksListStartFlag, "start", "s", time.Now().AddDate(0, 0, -8).Format(time.RFC3339), "Start")
	candlesticksListCmd.Flags().StringVarP(
		&candlesticksListEndFlag, "end", "E", time.Now().AddDate(0, 0, -1).Format(time.RFC3339), "End")
	candlesticksCmd.AddCommand(candlesticksListCmd)
	rootCmd.AddCommand(candlesticksCmd)
}
