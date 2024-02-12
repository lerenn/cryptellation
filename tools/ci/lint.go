package main

import (
	"context"

	"dagger.io/dagger"
	"github.com/lerenn/cryptellation/pkg/ci"
	"github.com/spf13/cobra"
)

func linters() map[string]*dagger.Container {
	return map[string]*dagger.Container{
		"cmd/cryptellation":     ci.Linter(client, "/cmd/cryptellation"),
		"cmd/cryptellation-tui": ci.Linter(client, "/cmd/cryptellation-tui"),
		"pkg":                   ci.Linter(client, "/pkg"),
		"svc/backtests":         ci.Linter(client, "/svc/backtests"),
		"svc/candlesticks":      ci.Linter(client, "/svc/candlesticks"),
		"svc/exchanges":         ci.Linter(client, "/svc/exchanges"),
		"svc/indicators":        ci.Linter(client, "/svc/indicators"),
		"svc/ticks":             ci.Linter(client, "/svc/ticks"),
		"tools/ci":              ci.Linter(client, "/tools/ci"),
		"tools/tag":             ci.Linter(client, "/tools/tag"),
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
