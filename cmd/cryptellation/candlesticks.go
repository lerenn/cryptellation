package main

import (
	"context"
	"fmt"
	"time"

	client "github.com/lerenn/cryptellation/clients/go"
	"github.com/lerenn/cryptellation/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/pkg/models/period"
	"github.com/lerenn/cryptellation/pkg/utils"
	"github.com/spf13/cobra"
)

var candlesticksCmd = &cobra.Command{
	Use:     "candlesticks",
	Aliases: []string{"c"},
	Short:   "Manipulate candlesticks service",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return executeParentPersistentPreRuns(cmd, args)
	},
}

var candlesticksReadCmd = &cobra.Command{
	Use:     "read",
	Aliases: []string{"r"},
	Short:   "Read candlesticks from service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		list, err := services.Candlesticks.Read(context.Background(), client.ReadCandlesticksPayload{
			ExchangeName: "binance",
			PairSymbol:   "ETH-USDT",
			Period:       period.H1,
			Start:        utils.ToReference(time.Now().AddDate(0, 0, -7)),
			End:          utils.ToReference(time.Now()),
		})
		if err != nil {
			return err
		}

		return list.Loop(func(t time.Time, cs candlestick.Candlestick) (bool, error) {
			fmt.Println("", t, "||", cs.String())
			return false, nil
		})
	},
}

func initCandlesticks(rootCmd *cobra.Command) {
	candlesticksCmd.AddCommand(candlesticksReadCmd)
	rootCmd.AddCommand(candlesticksCmd)
}
