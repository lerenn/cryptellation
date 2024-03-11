package main

import (
	"context"
	"fmt"

	"dagger.io/dagger"
	"github.com/lerenn/cryptellation/pkg/ci"
	backtestsCi "github.com/lerenn/cryptellation/svc/backtests/pkg/ci"
	candlesticksCi "github.com/lerenn/cryptellation/svc/candlesticks/pkg/ci"
	exchangesCi "github.com/lerenn/cryptellation/svc/exchanges/pkg/ci"
	indicatorsCi "github.com/lerenn/cryptellation/svc/indicators/pkg/ci"
	ticksCi "github.com/lerenn/cryptellation/svc/ticks/pkg/ci"
	"github.com/spf13/cobra"
)

var (
	testsTypeFlag string
)

func unitTests() map[string]*dagger.Container {
	return map[string]*dagger.Container{
		pathCmdCli:          ci.UnitTests(client, pathCmdCli),
		pathCmdTui:          ci.UnitTests(client, pathCmdTui),
		pathPkg:             ci.UnitTests(client, pathPkg),
		pathSvcBacktests:    backtestsCi.UnitTests(client),
		pathSvcCandlesticks: candlesticksCi.UnitTests(client),
		pathSvcExchanges:    exchangesCi.UnitTests(client),
		pathSvcIndicators:   indicatorsCi.UnitTests(client),
		pathSvcTicks:        ticksCi.UnitTests(client),
		pathToolsCi:         ci.UnitTests(client, pathToolsCi),
	}
}

func integrationTests() map[string]*dagger.Container {
	return map[string]*dagger.Container{
		pathSvcBacktests:    backtestsCi.IntegrationTests(client),
		pathSvcCandlesticks: candlesticksCi.IntegrationTests(client),
		pathSvcExchanges:    exchangesCi.IntegrationTests(client),
		pathSvcIndicators:   indicatorsCi.IntegrationTests(client),
		pathSvcTicks:        ticksCi.IntegrationTests(client),
	}
}

func endToEndTests() map[string]*dagger.Container {
	return map[string]*dagger.Container{
		pathSvcBacktests:    backtestsCi.EndToEndTests(client),
		pathSvcCandlesticks: candlesticksCi.EndToEndTests(client),
		pathSvcExchanges:    exchangesCi.EndToEndTests(client),
		pathSvcIndicators:   indicatorsCi.EndToEndTests(client),
		pathSvcTicks:        ticksCi.EndToEndTests(client),
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
