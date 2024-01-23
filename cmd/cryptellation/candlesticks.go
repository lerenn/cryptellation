package main

import (
	"context"
	"fmt"
	"time"

	"github.com/lerenn/cryptellation/pkg/utils"
	client "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"
	candlesticks "github.com/lerenn/cryptellation/svc/candlesticks/clients/go/nats"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"
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

		candlesticksClient, err = candlesticks.NewClient(globalConfig)
		if err != nil {
			return fmt.Errorf("error when creating new candlesticks client: %w", err)
		}

		return nil
	},
}

var candlesticksReadCmd = &cobra.Command{
	Use:     "read",
	Aliases: []string{"r"},
	Short:   "Read candlesticks from service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		list, err := candlesticksClient.Read(context.Background(), client.ReadCandlesticksPayload{
			ExchangeName: "binance",
			PairSymbol:   "ETH-USDT",
			Period:       period.H1,
			Start:        utils.ToReference(time.Now().AddDate(0, 0, -8)),
			End:          utils.ToReference(time.Now().AddDate(0, 0, -1)),
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
		info, err := candlesticksClient.ServiceInfo(context.TODO())
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
