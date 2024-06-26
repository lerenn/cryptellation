package main

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var forwardtestsCmd = &cobra.Command{
	Use:     "forwardtests",
	Aliases: []string{"c"},
	Short:   "Manipulate forwardtests service",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if err := executeParentPersistentPreRuns(cmd, args); err != nil {
			return err
		}

		return nil
	},
}

var forwardtestsInfoCmd = &cobra.Command{
	Use:     "service",
	Aliases: []string{"service"},
	Short:   "Read info from forwardtests service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return displayServiceInfo(globalClient.ForwardTests())
	},
}

var forwardtestsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"list"},
	Short:   "List forward tests",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		list, err := globalClient.ForwardTests().ListForwardTests(cmd.Context())
		if err != nil {
			return err
		}

		for _, l := range list {
			cmd.Println(l)
		}
		return nil
	},
}

var forwardtestsStatusCmd = &cobra.Command{
	Use:     "status <id>",
	Aliases: []string{"status"},
	Args:    cobra.ExactArgs(1),
	Short:   "Get forward test status",
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := uuid.Parse(args[0])
		if err != nil {
			return err
		}

		status, err := globalClient.ForwardTests().GetStatus(cmd.Context(), id)
		if err != nil {
			return err
		}

		fmt.Printf("%+v\n", status)
		return nil
	},
}

func initForwardTests(rootCmd *cobra.Command) {
	forwardtestsCmd.AddCommand(forwardtestsListCmd)
	forwardtestsCmd.AddCommand(forwardtestsInfoCmd)
	forwardtestsCmd.AddCommand(forwardtestsStatusCmd)
	rootCmd.AddCommand(forwardtestsCmd)
}
