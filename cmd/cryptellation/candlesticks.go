package main

import (
	"context"
	"fmt"
	"time"

	"github.com/lerenn/cryptellation/pkg/utils"
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
	Use:     "info",
	Aliases: []string{"info"},
	Short:   "Read info from candlesticks service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		info, err := services.Candlesticks.ServiceInfo(context.TODO())
		if err != nil {
			return err
		}

		fmt.Printf("%+v\n", info)
		return nil
	},
}

func initCandlesticks(rootCmd *cobra.Command) {
	candlesticksCmd.AddCommand(candlesticksReadCmd)
	candlesticksCmd.AddCommand(candlesticksInfoCmd)
	rootCmd.AddCommand(candlesticksCmd)
}
