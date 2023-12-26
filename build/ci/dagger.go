package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
	"github.com/lerenn/cryptellation/pkg/pipeline"
	"github.com/spf13/cobra"
)

var (
	testsTypeFlag string
)

var (
	client *dagger.Client

	generator *dagger.Container
	linter    *dagger.Container

	unitTests                               []*dagger.Container
	integrationPreFlights, integrationTests []*dagger.Container
	endToEndPreFlights, endToEndTests       []*dagger.Container
)

var rootCmd = &cobra.Command{
	Use:   "./build/ci/dagger.go",
	Short: "A simple CLI to execute cryptellation project CI/CD with dagger",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		ctx := context.TODO()

		// Initialize Dagger client
		client, err = dagger.Connect(context.Background(), dagger.WithLogOutput(os.Stderr))
		if err != nil {
			return err
		}
		defer client.Close()

		// Create containers
		generator = pipeline.Generator(client)
		linter = pipeline.Linter(client)
		unitTests = pipeline.UnitTests(client)
		integrationPreFlights, integrationTests = pipeline.IntegrationTests(ctx, client)
		endToEndPreFlights, endToEndTests = pipeline.EndToEndTests(ctx, client)

		return nil
	},
}

var allCmd = &cobra.Command{
	Use:     "all",
	Aliases: []string{"a"},
	Short:   "Execute all CI steps",
	Run: func(cmd *cobra.Command, args []string) {
		pipeline.ExecuteContainersInParallel(context.Background(), []*dagger.Container{generator, linter})
		pipeline.ExecuteContainersInParallel(context.Background(), unitTests, integrationTests)
	},
}

var generatorCmd = &cobra.Command{
	Use:     "generator",
	Aliases: []string{"g"},
	Short:   "Execute generator step of the CI",
	Run: func(cmd *cobra.Command, args []string) {
		pipeline.ExecuteContainersInParallel(context.Background(), []*dagger.Container{generator})
	},
}

var linterCmd = &cobra.Command{
	Use:     "linter",
	Aliases: []string{"g"},
	Short:   "Execute linter step of the CI",
	Run: func(cmd *cobra.Command, args []string) {
		pipeline.ExecuteContainersInParallel(context.Background(), []*dagger.Container{linter})
	},
}

var testCmd = &cobra.Command{
	Use:     "test",
	Aliases: []string{"g"},
	Short:   "Execute tests step of the CI",
	Run: func(cmd *cobra.Command, args []string) {
		switch testsTypeFlag {
		case "unit":
			pipeline.ExecuteContainersInParallel(context.Background(), unitTests)
		case "integration":
			pipeline.ExecuteContainersInParallel(context.Background(), integrationPreFlights)
			pipeline.ExecuteContainersInParallel(context.Background(), integrationTests)
		case "end-to-end":
			pipeline.ExecuteContainersInParallel(context.Background(), endToEndPreFlights)
			pipeline.ExecuteContainersInParallel(context.Background(), endToEndTests)
		case "all:":
			pipeline.ExecuteContainersInParallel(context.Background(), endToEndPreFlights, integrationPreFlights)
			pipeline.ExecuteContainersInParallel(context.Background(), unitTests, integrationTests, endToEndTests)
		default:
			panic(fmt.Sprintf("invalid type %q", testsTypeFlag))
		}
	},
}

func main() {
	rootCmd.AddCommand(allCmd)
	rootCmd.AddCommand(generatorCmd)
	rootCmd.AddCommand(linterCmd)
	rootCmd.AddCommand(testCmd)
	testCmd.Flags().StringVarP(&testsTypeFlag, "type", "t", "all", "Type of the test (all, unit, integration, end-to-end)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
