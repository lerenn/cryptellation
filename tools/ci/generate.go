package main

import (
	"context"

	"github.com/lerenn/cryptellation/pkg/ci"
	"github.com/spf13/cobra"
)

func generators() map[string]func(context.Context) error {
	return map[string]func(context.Context) error{
		"svc/backtests":    ci.Generator(client, "svc/backtests"),
		"svc/candlesticks": ci.Generator(client, "svc/candlesticks"),
		"svc/exchanges":    ci.Generator(client, "svc/exchanges"),
		"svc/indicators":   ci.Generator(client, "svc/indicators"),
		"svc/ticks":        ci.Generator(client, "svc/ticks"),
	}
}

func runGenerators(cmd *cobra.Command, args []string) {
	ci.ExecuteInParallel(context.Background(), filterWithPath(generators())...)
}

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"g"},
	Short:   "Execute generator step of the CI",
	Run:     runGenerators,
}

func addGenerateCmdTo(cmd *cobra.Command) {
	cmd.AddCommand(generateCmd)
}
