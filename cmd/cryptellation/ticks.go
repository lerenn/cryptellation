package main

import (
	"context"
	"fmt"
	"time"

	"cryptellation/pkg/models/event"

	"github.com/spf13/cobra"
)

var (
	ticksExchange string
	ticksPair     string
	ticksTimeout  string
)

var ticksCmd = &cobra.Command{
	Use:     "ticks",
	Aliases: []string{"t"},
	Short:   "Manipulate ticks service",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if err := executeParentPersistentPreRuns(cmd, args); err != nil {
			return err
		}

		return nil
	},
}

var ticksListenCmd = &cobra.Command{
	Use:     "listen",
	Aliases: []string{"l"},
	Short:   "Listen to ticks on service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		timeout, err := time.ParseDuration(ticksTimeout)
		if err != nil {
			return err
		}

		ts := event.TickSubscription{
			Exchange: ticksExchange,
			Pair:     ticksPair,
		}
		ch, err := globalClient.Ticks().SubscribeToTicks(context.TODO(), ts)
		if err != nil {
			return err
		}

		for {
			select {
			case t := <-ch:
				fmt.Println(t.String())
			case <-time.After(timeout):
				return fmt.Errorf("timeout after %s", timeout)
			}
		}
	},
}

var ticksInfoCmd = &cobra.Command{
	Use:     "service",
	Aliases: []string{"svc"},
	Short:   "Read info from ticks service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return displayServiceInfo(globalClient.Ticks())
	},
}

func initTicks(rootCmd *cobra.Command) {
	ticksListenCmd.Flags().StringVarP(&ticksExchange, "exchange", "e", "binance", "Exchange to listen to")
	ticksListenCmd.Flags().StringVarP(&ticksPair, "pair", "p", "BTC-USDT", "Pair to listen to")
	ticksListenCmd.Flags().StringVarP(&ticksTimeout, "timeout", "t", "30s", "Timeout for listening to ticks")
	ticksCmd.AddCommand(ticksListenCmd)

	ticksCmd.AddCommand(ticksInfoCmd)

	rootCmd.AddCommand(ticksCmd)
}
