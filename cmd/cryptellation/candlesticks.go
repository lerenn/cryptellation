package main

import (
	"context"
	"fmt"
	"time"

	client "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"

	"github.com/spf13/cobra"
)

var candlesticksCmd = &cobra.Command{
	Use:     "candlesticks",
	Aliases: []string{"c"},
	Short:   "Manipulate candlesticks service",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if err := executeParentPersistentPreRuns(cmd, args); err != nil {
			return err
		}

		return nil
	},
}

var (
	candlesticksReadExchangeFlag string
	candlesticksReadPairFlag     string
	candlesticksReadPeriodFlag   string
	candlesticksReadStartFlag    string
	candlesticksReadEndFlag      string
)

var candlesticksReadCmd = &cobra.Command{
	Use:     "read",
	Aliases: []string{"r"},
	Short:   "Read candlesticks from service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		start, err := time.Parse(time.RFC3339, candlesticksReadStartFlag)
		if err != nil {
			return err
		}

		end, err := time.Parse(time.RFC3339, candlesticksReadEndFlag)
		if err != nil {
			return err
		}

		list, err := globalClient.Candlesticks.Read(context.Background(), client.ReadCandlesticksPayload{
			Exchange: candlesticksReadExchangeFlag,
			Pair:     candlesticksReadPairFlag,
			Period:   period.Symbol(candlesticksReadPeriodFlag),
			Start:    &start,
			End:      &end,
		})
		if err != nil {
			return err
		}

		switch {
		case jsonOutput:
			return displayJSON(list.ToArray())
		default:
			if err := list.Loop(func(cs candlestick.Candlestick) (bool, error) {
				fmt.Println(cs.String())
				return false, nil
			}); err != nil {
				return err
			}
		}

		return nil
	},
}

var candlesticksInfoCmd = &cobra.Command{
	Use:     "service",
	Aliases: []string{"svc"},
	Short:   "Read info from candlesticks service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return displayServiceInfo(globalClient.Candlesticks)
	},
}

func initCandlesticks(rootCmd *cobra.Command) {
	candlesticksReadCmd.Flags().StringVarP(&candlesticksReadExchangeFlag, "exchange", "e", "binance", "Exchange")
	candlesticksReadCmd.Flags().StringVarP(&candlesticksReadPairFlag, "pair", "p", "ETH-USDT", "Pair")
	candlesticksReadCmd.Flags().StringVarP(&candlesticksReadPeriodFlag, "period", "P", "H1", "Period")
	candlesticksReadCmd.Flags().StringVarP(&candlesticksReadStartFlag, "start", "s", time.Now().AddDate(0, 0, -8).Format(time.RFC3339), "Start")
	candlesticksReadCmd.Flags().StringVarP(&candlesticksReadEndFlag, "end", "E", time.Now().AddDate(0, 0, -1).Format(time.RFC3339), "End")
	candlesticksCmd.AddCommand(candlesticksReadCmd)

	candlesticksCmd.AddCommand(candlesticksInfoCmd)
	rootCmd.AddCommand(candlesticksCmd)
}
