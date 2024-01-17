package main

import (
	"context"

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

func linters() map[string]*dagger.Container {
	return map[string]*dagger.Container{
		"clients":          clientsCi.Linter(client),
		"cmd":              cmdCi.Linter(client),
		"pkg":              pkgCi.Linter(client),
		"svc/backtests":    backtestsCi.Linter(client),
		"svc/candlesticks": candlesticksCi.Linter(client),
		"svc/exchanges":    exchangesCi.Linter(client),
		"svc/indicators":   indicatorsCi.Linter(client),
		"svc/ticks":        ticksCi.Linter(client),
	}
}

func runLinters(cmd *cobra.Command, args []string) {
	ci.ExecuteContainersInParallel(
		context.Background(),
		filterWithPath(linters()),
	)
}

var lintCmd = &cobra.Command{
	Use:     "lint",
	Aliases: []string{"l"},
	Short:   "Execute linter step of the CI",
	Run:     runLinters,
}

func addLintCmdTo(cmd *cobra.Command) {
	cmd.AddCommand(lintCmd)
}
