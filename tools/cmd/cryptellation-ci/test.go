package main

import (
	"context"
	"fmt"

	"dagger.io/dagger"
	backtestsCi "github.com/lerenn/cryptellation/svc/backtests/pkg/ci"
	candlesticksCi "github.com/lerenn/cryptellation/svc/candlesticks/pkg/ci"
	exchangesCi "github.com/lerenn/cryptellation/svc/exchanges/pkg/ci"
	indicatorsCi "github.com/lerenn/cryptellation/svc/indicators/pkg/ci"
	ticksCi "github.com/lerenn/cryptellation/svc/ticks/pkg/ci"
	"github.com/lerenn/cryptellation/tools/pkg/ci"
	"github.com/spf13/cobra"
)

var (
	testsTypeFlag string
)

func unitTests() map[string]*dagger.Container {
	return map[string]*dagger.Container{
		"cmd/cryptellation":     ci.UnitTests(client, "/cmd/cryptellation"),
		"cmd/cryptellation-tui": ci.UnitTests(client, "/cmd/cryptellation-tui"),
		"pkg":                   ci.UnitTests(client, "/pkg"),
		"svc/backtests":         backtestsCi.UnitTests(client),
		"svc/candlesticks":      candlesticksCi.UnitTests(client),
		"svc/exchanges":         exchangesCi.UnitTests(client),
		"svc/indicators":        indicatorsCi.UnitTests(client),
		"svc/ticks":             ticksCi.UnitTests(client),
	}
}

func integrationTests() map[string]*dagger.Container {
	return map[string]*dagger.Container{
		"svc/backtests":    backtestsCi.IntegrationTests(client),
		"svc/candlesticks": candlesticksCi.IntegrationTests(client),
		"svc/exchanges":    exchangesCi.IntegrationTests(client),
		"svc/indicators":   indicatorsCi.IntegrationTests(client),
		"svc/ticks":        ticksCi.IntegrationTests(client),
	}
}

func endToEndTests() map[string]*dagger.Container {
	return map[string]*dagger.Container{
		"svc/backtests":    backtestsCi.EndToEndTests(client),
		"svc/candlesticks": candlesticksCi.EndToEndTests(client),
		"svc/exchanges":    exchangesCi.EndToEndTests(client),
		"svc/indicators":   indicatorsCi.EndToEndTests(client),
		"svc/ticks":        ticksCi.EndToEndTests(client),
	}
}

func runTests(cmd *cobra.Command, args []string) error {
	ut := unitTests()
	it := integrationTests()
	et := endToEndTests()

	switch testsTypeFlag {
	case "unit":
		ci.ExecuteContainersInParallel(context.Background(), filterWithPath(ut))
	case "integration":
		ci.ExecuteContainersInParallel(context.Background(), filterWithPath(it))
	case "end-to-end":
		ci.ExecuteContainersInParallel(context.Background(), filterWithPath(et))
	case "all":
		ci.ExecuteContainersInParallel(context.Background(), filterWithPath(ut))
		ci.ExecuteContainersInParallel(context.Background(), filterWithPath(it))
		ci.ExecuteContainersInParallel(context.Background(), filterWithPath(et))
	default:
		return fmt.Errorf("invalid type %q", testsTypeFlag)
	}

	return nil
}

var testCmd = &cobra.Command{
	Use:     "test",
	Aliases: []string{"t"},
	Short:   "Execute tests step of the CI",
	RunE:    runTests,
}

func addTestCmdTo(cmd *cobra.Command) {
	testCmd.Flags().StringVarP(&testsTypeFlag, "type", "t", "all", "Type of the test (all, unit, integration, end-to-end)")
	cmd.AddCommand(testCmd)
}
