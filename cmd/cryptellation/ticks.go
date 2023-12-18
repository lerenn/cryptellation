package main

import (
	"context"
	"fmt"
	"time"

	client "github.com/lerenn/cryptellation/clients/go"
	"github.com/spf13/cobra"
)

var ticksCmd = &cobra.Command{
	Use:     "ticks",
	Aliases: []string{"c"},
	Short:   "Manipulate ticks service",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return executeParentPersistentPreRuns(cmd, args)
	},
}

var ticksRegisterCmd = &cobra.Command{
	Use:     "register",
	Aliases: []string{"r"},
	Short:   "Register to ticks on service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		err = services.Ticks.Register(context.Background(), client.TicksFilterPayload{
			ExchangeName: "binance",
			PairSymbol:   "BTC-USDT",
		})

		return err
	},
}

var WatchTicksCmd = &cobra.Command{
	Use:     "listen",
	Aliases: []string{"r"},
	Short:   "Listen to ticks on service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		ch, err := services.Ticks.Listen(context.Background(), client.TicksFilterPayload{
			ExchangeName: "binance",
			PairSymbol:   "BTC-USDT",
		})
		if err != nil {
			return err
		}

		for {
			select {
			case t := <-ch:
				fmt.Println(t.String())
			case <-time.After(10 * time.Second):
				return fmt.Errorf("Timeout")
			}
		}
	},
}

var UnregisterToTicksRequestCmd = &cobra.Command{
	Use:     "unregister",
	Aliases: []string{"r"},
	Short:   "Unregister to ticks on service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		err = services.Ticks.Unregister(context.Background(), client.TicksFilterPayload{
			ExchangeName: "binance",
			PairSymbol:   "BTC-USDT",
		})

		return err
	},
}

func initTicks(rootCmd *cobra.Command) {
	ticksCmd.AddCommand(ticksRegisterCmd)
	ticksCmd.AddCommand(WatchTicksCmd)
	ticksCmd.AddCommand(UnregisterToTicksRequestCmd)
	rootCmd.AddCommand(ticksCmd)
}
