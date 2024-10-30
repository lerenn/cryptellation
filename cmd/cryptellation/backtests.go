package main

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var backtestsCmd = &cobra.Command{
	Use:     "backtests",
	Aliases: []string{"bt"},
	Short:   "Manipulate backtests service",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if err := executeParentPersistentPreRuns(cmd, args); err != nil {
			return err
		}

		return nil
	},
}

var backtestsList = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "ls"},
	Short:   "List backtests",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		list, err := globalClient.Backtests.List(cmd.Context())
		if err != nil {
			return err
		}

		switch {
		case jsonOutput:
			return displayJSON(list)
		default:
			fmt.Printf("%-40s\n", "ID")
			for _, bt := range list {
				fmt.Printf("%-40s\n", bt.ID)
			}
		}
		return nil
	},
}

var backtestsGetCmd = &cobra.Command{
	Use:     "get <id>",
	Aliases: []string{"s"},
	Args:    cobra.ExactArgs(1),
	Short:   "Get backtest",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		id, err := uuid.Parse(args[0])
		if err != nil {
			return err
		}

		bt, err := globalClient.Backtests.Get(cmd.Context(), id)
		if err != nil {
			return err
		}

		switch {
		case jsonOutput:
			return displayJSON(bt)
		default:
			fmt.Printf("ID:\t\t%s\n", bt.ID)
			fmt.Printf("Start:\t\t%s\n", bt.Parameters.StartTime)
			fmt.Printf("End:\t\t%s\n", bt.Parameters.EndTime)
			fmt.Printf("Period:\t\t%s\n", bt.Parameters.PricePeriod)
			fmt.Printf("Tick subs:\n")
			for _, ts := range bt.PricesSubscriptions {
				fmt.Printf("\t%s: %s\n", ts.Exchange, ts.Pair)
			}
		}

		return nil
	},
}

var backtestsOrdersCmd = &cobra.Command{
	Use:     "orders",
	Aliases: []string{"o"},
	Short:   "Manage backtest orders",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if err := executeParentPersistentPreRuns(cmd, args); err != nil {
			return err
		}

		return nil
	},
}

var backtestsOrdersListCmd = &cobra.Command{
	Use:     "list <id>",
	Aliases: []string{"l", "ls"},
	Args:    cobra.ExactArgs(1),
	Short:   "List backtest orders",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		id, err := uuid.Parse(args[0])
		if err != nil {
			return err
		}

		orders, err := globalClient.Backtests.ListOrders(cmd.Context(), id)
		if err != nil {
			return err
		}

		switch {
		case jsonOutput:
			return displayJSON(orders)
		default:
			fmt.Printf("%-40s\t%-40s\t%-40s\t%-40s\n", "ID", "Time", "Type", "Quantity")
			for _, o := range orders {
				fmt.Printf("%-40s\t%-40s\t%-40s\t%-40f\n", o.ID, o.ExecutionTime, o.Type, o.Quantity)
			}
		}

		return nil
	},
}

var backtestsServiceInfoCmd = &cobra.Command{
	Use:     "service",
	Aliases: []string{"svc"},
	Short:   "Read info from backtests service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return displayServiceInfo(globalClient.Backtests)
	},
}

func initBacktests(rootCmd *cobra.Command) {
	backtestsCmd.AddCommand(backtestsGetCmd)
	backtestsCmd.AddCommand(backtestsList)
	backtestsCmd.AddCommand(backtestsServiceInfoCmd)

	backtestsOrdersCmd.AddCommand(backtestsOrdersListCmd)
	backtestsCmd.AddCommand(backtestsOrdersCmd)

	rootCmd.AddCommand(backtestsCmd)
}
