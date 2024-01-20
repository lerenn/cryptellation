package main

import (
	"context"
	"fmt"
	"time"

	client "github.com/lerenn/cryptellation/svc/ticks/clients/go"
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

var ticksWatchCmd = &cobra.Command{
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

var ticksUnregisterCmd = &cobra.Command{
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

var ticksInfoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"info"},
	Short:   "Read info from ticks service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		info, err := services.Ticks.ServiceInfo(context.TODO())
		if err != nil {
			return err
		}

		fmt.Printf("%+v\n", info)
		return nil
	},
}

func initTicks(rootCmd *cobra.Command) {
	ticksCmd.AddCommand(ticksInfoCmd)
	ticksCmd.AddCommand(ticksRegisterCmd)
	ticksCmd.AddCommand(ticksUnregisterCmd)
	ticksCmd.AddCommand(ticksWatchCmd)

	rootCmd.AddCommand(ticksCmd)
}
