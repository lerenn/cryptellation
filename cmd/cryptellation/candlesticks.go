package main

import (
	"context"
	"fmt"
	"time"

	candlesticks "github.com/digital-feather/cryptellation/internal/candlesticks/ctrl/nats"
	"github.com/digital-feather/cryptellation/pkg/candlestick"
	"github.com/digital-feather/cryptellation/pkg/period"
	"github.com/digital-feather/cryptellation/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	candlesticksClient candlesticks.Client
)

var candlesticksCmd = &cobra.Command{
	Use:     "candlesticks",
	Aliases: []string{"c"},
	Short:   "Manipulate candlesticks service",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if err := executeParentPersistentPreRuns(cmd, args); err != nil {
			return err
		}

		candlesticksClient, err = candlesticks.New(globalNATSConfig)
		return err
	},
}

var candlesticksReadCmd = &cobra.Command{
	Use:     "read",
	Aliases: []string{"r"},
	Short:   "Read candlesticks from service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		list, err := candlesticksClient.ReadCandlesticks(context.Background(), candlesticks.ReadCandlesticksPayload{
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
