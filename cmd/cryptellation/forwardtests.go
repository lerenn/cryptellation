package main

import "github.com/spf13/cobra"

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
	Use:     "info",
	Aliases: []string{"info"},
	Short:   "Read info from forwardtests service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return displayServiceInfo(globalClient.ForwardTests())
	},
}

func initForwardTests(rootCmd *cobra.Command) {
	forwardtestsCmd.AddCommand(forwardtestsInfoCmd)
	rootCmd.AddCommand(forwardtestsCmd)
}
