package main

import (
	"context"
	"fmt"
	"time"

	"cryptellation/pkg/utils"

	client "cryptellation/svc/candlesticks/clients/go"
	"cryptellation/svc/candlesticks/pkg/candlestick"
	"cryptellation/svc/candlesticks/pkg/period"

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

var candlesticksReadCmd = &cobra.Command{
	Use:     "read",
	Aliases: []string{"r"},
	Short:   "Read candlesticks from service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		list, err := globalClient.Candlesticks().Read(context.Background(), client.ReadCandlesticksPayload{
			Exchange: "binance",
			Pair:     "ETH-USDT",
			Period:   period.H1,
			Start:    utils.ToReference(time.Now().AddDate(0, 0, -8)),
			End:      utils.ToReference(time.Now().AddDate(0, 0, -1)),
		})
		if err != nil {
			return err
		}

		if err := list.Loop(func(t time.Time, cs candlestick.Candlestick) (bool, error) {
			fmt.Println("", t, "||", cs.String())
			return false, nil
		}); err != nil {
			return err
		}

		return nil
	},
}

var candlesticksInfoCmd = &cobra.Command{
	Use:     "service",
	Aliases: []string{"svc"},
	Short:   "Read info from candlesticks service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return displayServiceInfo(globalClient.Candlesticks())
	},
}

func initCandlesticks(rootCmd *cobra.Command) {
	candlesticksCmd.AddCommand(candlesticksReadCmd)
	candlesticksCmd.AddCommand(candlesticksInfoCmd)
	rootCmd.AddCommand(candlesticksCmd)
}
