package main

import (
	"context"
	"fmt"

	"dagger.io/dagger"
	clientsCi "github.com/lerenn/cryptellation/cmd/ci/internal/clients"
	cmdCi "github.com/lerenn/cryptellation/cmd/ci/internal/cmd"
	pkgCi "github.com/lerenn/cryptellation/cmd/ci/internal/pkg"
	"github.com/lerenn/cryptellation/pkg/ci"
	backtestsCi "github.com/lerenn/cryptellation/svc/backtests/pkg/ci"
	candlesticksCi "github.com/lerenn/cryptellation/svc/candlesticks/pkg/ci"
	exchangesCi "github.com/lerenn/cryptellation/svc/exchanges/pkg/ci"
	indicatorsCi "github.com/lerenn/cryptellation/svc/indicators/pkg/ci"
	ticksCi "github.com/lerenn/cryptellation/svc/ticks/pkg/ci"
	"github.com/spf13/cobra"
)

func unitTests() map[string]*dagger.Container {
	return map[string]*dagger.Container{
		"clients":          clientsCi.UnitTests(client),
		"cmd":              cmdCi.UnitTests(client),
		"pkg":              pkgCi.UnitTests(client),
		"svc/backtests":    backtestsCi.UnitTests(client),
		"svc/candlesticks": candlesticksCi.UnitTests(client),
		"svc/exchanges":    exchangesCi.UnitTests(client),
		"svc/indicators":   indicatorsCi.UnitTests(client),
		"svc/ticks":        ticksCi.UnitTests(client),
	}
}

func integrationTests() map[string]*dagger.Container {
	return map[string]*dagger.Container{
		"pkg":              pkgCi.IntegrationTests(client),
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
		ci.ExecuteContainersInParallel(context.Background(), filterContainerWithPath(ut))
	case "integration":
		ci.ExecuteContainersInParallel(context.Background(), filterContainerWithPath(it))
	case "end-to-end":
		ci.ExecuteContainersInParallel(context.Background(), filterContainerWithPath(et))
	case "all":
		ci.ExecuteContainersInParallel(context.Background(), filterContainerWithPath(ut))
		ci.ExecuteContainersInParallel(context.Background(), filterContainerWithPath(it))
		ci.ExecuteContainersInParallel(context.Background(), filterContainerWithPath(et))
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
