package main

import (
	"context"
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:     "status",
	Aliases: []string{"status"},
	Short:   "Read info from services",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		// Get services info
		infos, err := globalClient.ServicesInfo(context.TODO())
		if err != nil {
			return err
		}

		// Get keys in alphabetical order
		mk := make([]string, 0, len(infos))
		for k := range infos {
			mk = append(mk, k)
		}
		sort.Strings(mk)

		// Print info in alphabetical order
		for _, k := range mk {
			fmt.Printf("%s: %+v\n", k, infos[k])
		}

		return nil
	},
}

func initStatuses(rootCmd *cobra.Command) {
	rootCmd.AddCommand(statusCmd)
}
