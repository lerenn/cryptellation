package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var exchangesCmd = &cobra.Command{
	Use:     "exchanges",
	Aliases: []string{"c"},
	Short:   "Manipulate exchanges service",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return executeParentPersistentPreRuns(cmd, args)
	},
}

var exchangesReadCmd = &cobra.Command{
	Use:     "read",
	Aliases: []string{"r"},
	Short:   "Read exchanges from service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		list, err := services.Exchanges.Read(context.Background(), "binance")
		if err != nil {
			return err
		}

		fmt.Println(list)
		return nil
	},
}

var exchangesInfoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"info"},
	Short:   "Read info from exchanges service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		info, err := services.Exchanges.ServiceInfo(context.TODO())
		if err != nil {
			return err
		}

		fmt.Printf("%+v\n", info)
		return nil
	},
}

func initExchanges(rootCmd *cobra.Command) {
	exchangesCmd.AddCommand(exchangesReadCmd)
	exchangesCmd.AddCommand(exchangesInfoCmd)
	rootCmd.AddCommand(exchangesCmd)
}
