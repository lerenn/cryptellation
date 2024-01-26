package main

import (
	"context"

	"github.com/lerenn/cryptellation/pkg/ci"
	"github.com/spf13/cobra"
)

func updators() map[string]func(context.Context) error {
	return map[string]func(context.Context) error{
		"cmd/cryptellation":     ci.UpdateGoMod(client, "cmd/cryptellation"),
		"cmd/cryptellation-tui": ci.UpdateGoMod(client, "cmd/cryptellation-tui"),
		"pkg":                   ci.UpdateGoMod(client, "pkg"),
		"svc/backtests":         ci.UpdateGoMod(client, "svc/backtests"),
		"svc/candlesticks":      ci.UpdateGoMod(client, "svc/candlesticks"),
		"svc/exchanges":         ci.UpdateGoMod(client, "svc/exchanges"),
		"svc/indicators":        ci.UpdateGoMod(client, "svc/indicators"),
		"svc/ticks":             ci.UpdateGoMod(client, "svc/ticks"),
		"tools/ci":              ci.UpdateGoMod(client, "tools/ci"),
		"tools/tag":             ci.UpdateGoMod(client, "tools/tag"),
	}
}

func runUpdators(cmd *cobra.Command, args []string) {
	ci.ExecuteInParallel(context.Background(), filterWithPath(updators())...)
}

var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"u"},
	Short:   "Execute updates step of the CI",
	Run:     runUpdators,
}

func addUpdateCmdTo(cmd *cobra.Command) {
	cmd.AddCommand(updateCmd)
}
