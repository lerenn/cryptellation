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

var exchangesListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List exchanges",
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := cryptellationClient.ListExchanges(cmd.Context(), api.ListExchangesParams{
			Names: args,
		})
		if err != nil {
			return err
		}

		switch {
		case jsonOutput:
			return displayJSON(cmd.Context(), res.List)
		default:
			for i := range res.List {
				fmt.Println(res.List[i].String())
			}
		}

		return nil
	},
}

func addExchangesCommands() {
	exchangesCmd.AddCommand(exchangesListCmd)
	rootCmd.AddCommand(exchangesCmd)
}
