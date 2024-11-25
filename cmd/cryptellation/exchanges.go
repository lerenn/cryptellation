package main

import (
	"fmt"

	"github.com/lerenn/cryptellation/v1/api"
	"github.com/spf13/cobra"
)

var exchangesCmd = &cobra.Command{
	Use:     "exchanges",
	Aliases: []string{"c"},
	Short:   "Manage exchanges",
}

var exchangesShowCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"s"},
	Short:   "Show exchange",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := cryptellationClient.GetExchange(cmd.Context(), api.GetExchangeParams{
			Name: args[0],
		})
		if err != nil {
			return err
		}

		switch {
		case jsonOutput:
			return displayJSON(cmd.Context(), res.Exchange)
		default:
			fmt.Println(res.Exchange.String())
		}

		return nil
	},
}

var exchangesListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List exchanges",
	RunE: func(cmd *cobra.Command, _ []string) error {
		res, err := cryptellationClient.ListExchanges(cmd.Context(), api.ListExchangesParams{})
		if err != nil {
			return err
		}

		switch {
		case jsonOutput:
			return displayJSON(cmd.Context(), res.List)
		default:
			fmt.Println("NAME")
			for i := range res.List {
				fmt.Println(res.List[i])
			}
		}

		return nil
	},
}

func addExchangesCommands() {
	exchangesCmd.AddCommand(exchangesShowCmd)
	exchangesCmd.AddCommand(exchangesListCmd)
	rootCmd.AddCommand(exchangesCmd)
}
