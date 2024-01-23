package main

import (
	"context"
	"fmt"
	"time"

	client "github.com/lerenn/cryptellation/svc/ticks/clients/go"
	ticks "github.com/lerenn/cryptellation/svc/ticks/clients/go/nats"
	"github.com/spf13/cobra"
)

var (
	ticksClient ticks.Client
)

var ticksCmd = &cobra.Command{
	Use:     "ticks",
	Aliases: []string{"c"},
	Short:   "Manipulate ticks service",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if err := executeParentPersistentPreRuns(cmd, args); err != nil {
			return err
		}

		ticksClient, err = ticks.NewClient(globalConfig)
		if err != nil {
			return fmt.Errorf("error when creating new candlesticks client: %w", err)
		}

		return nil
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

var ticksWatchCmd = &cobra.Command{
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

var ticksInfoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"info"},
	Short:   "Read info from ticks service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		info, err := ticksClient.ServiceInfo(context.TODO())
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
