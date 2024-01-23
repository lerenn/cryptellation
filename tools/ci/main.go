package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
	"github.com/spf13/cobra"
)

var (
	client *dagger.Client

	pathFlag string
)

var rootCmd = &cobra.Command{
	Use:   "dagger run go run ./tools/ci",
	Short: "A simple CLI to execute cryptellation project CI/CD with dagger",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		// Initialize Dagger client
		client, err = dagger.Connect(context.Background(), dagger.WithLogOutput(os.Stderr))
		if err != nil {
			return err
		}
		defer client.Close()

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		runGenerators(cmd, args)
		runLinters(cmd, args)
		runBuilders(cmd, args)

		if err := runTests(cmd, args); err != nil {
			return err
		}

		return nil
	},
}

func main() {
	addBuildCmdTo(rootCmd)
	addGenerateCmdTo(rootCmd)
	addLintCmdTo(rootCmd)
	addPublishCmdTo(rootCmd)
	addServeCmdTo(rootCmd)
	addTestCmdTo(rootCmd)
	addUpdateCmdTo(rootCmd)

	rootCmd.PersistentFlags().StringVarP(&pathFlag, "path", "p", "", "Specific part of the project to target")
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
