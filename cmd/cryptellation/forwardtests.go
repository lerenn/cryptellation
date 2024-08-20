package main

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var forwardtestsCmd = &cobra.Command{
	Use:     "forwardtests",
	Aliases: []string{"ft"},
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
	Aliases: []string{"svc"},
	Short:   "Read info from forwardtests service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return displayServiceInfo(globalClient.ForwardTests())
	},
}

var forwardtestsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List forward tests",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		list, err := globalClient.ForwardTests().ListForwardTests(cmd.Context())
		if err != nil {
			return err
		}

		fmt.Printf("%-40s%-20s\n", "ID", "UpdatedAt")
		for _, ft := range list {
			fmt.Printf("%-40s%-20s\n", ft.ID, ft.UpdatedAt)
		}
		return nil
	},
}

var forwardtestsStatusCmd = &cobra.Command{
	Use:     "status <id>",
	Aliases: []string{"s"},
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
