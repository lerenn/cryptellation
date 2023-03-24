package main

import (
	"context"
	"fmt"
	"time"

	client "github.com/digital-feather/cryptellation/clients/go"
	"github.com/digital-feather/cryptellation/clients/go/nats"
	"github.com/digital-feather/cryptellation/pkg/candlestick"
	"github.com/digital-feather/cryptellation/pkg/period"
	"github.com/digital-feather/cryptellation/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	candlesticks client.Candlesticks
)

var candlesticksCmd = &cobra.Command{
	Use:     "candlesticks",
	Aliases: []string{"c"},
	Short:   "Manipulate candlesticks service",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if err := executeParentPersistentPreRuns(cmd, args); err != nil {
			return err
		}

		candlesticks, err = nats.NewCandlesticks(globalNATSConfig)
		return err
	},
}

var candlesticksReadCmd = &cobra.Command{
	Use:     "read",
	Aliases: []string{"r"},
	Short:   "Read candlesticks from service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		list, err := candlesticks.Read(context.Background(), client.ReadCandlesticksPayload{
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
