package main

import (
	"context"
	"fmt"
	"time"

	client "github.com/digital-feather/cryptellation/clients/go"
	"github.com/digital-feather/cryptellation/clients/go/nats"
	"github.com/spf13/cobra"
)

var (
	ticksClient client.Ticks
)

var ticksCmd = &cobra.Command{
	Use:     "ticks",
	Aliases: []string{"c"},
	Short:   "Manipulate ticks service",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if err := executeParentPersistentPreRuns(cmd, args); err != nil {
			return err
		}

		ticksClient, err = nats.NewTicks(globalNATSConfig)
		return err
	},
}

var ticksRegisterCmd = &cobra.Command{
	Use:     "register",
	Aliases: []string{"r"},
	Short:   "Register to ticks on service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		err = ticksClient.Register(context.Background(), client.TicksFilterPayload{
			ExchangeName: "binance",
			PairSymbol:   "BTC-USDT",
		})

		return err
	},
}

var ticksListenCmd = &cobra.Command{
	Use:     "listen",
	Aliases: []string{"r"},
	Short:   "Listen to ticks on service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		ch, err := ticksClient.Listen(context.Background(), client.TicksFilterPayload{
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

var ticksUnregisterCmd = &cobra.Command{
	Use:     "unregister",
	Aliases: []string{"r"},
	Short:   "Unregister to ticks on service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		err = ticksClient.Unregister(context.Background(), client.TicksFilterPayload{
			ExchangeName: "binance",
			PairSymbol:   "BTC-USDT",
		})

		return err
	},
}

func initTicks(rootCmd *cobra.Command) {
	ticksCmd.AddCommand(ticksRegisterCmd)
	ticksCmd.AddCommand(ticksListenCmd)
	ticksCmd.AddCommand(ticksUnregisterCmd)
	rootCmd.AddCommand(ticksCmd)
}
