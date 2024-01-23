package main

import (
	"context"

	"github.com/lerenn/cryptellation/pkg/ci"
	backtestsCi "github.com/lerenn/cryptellation/svc/backtests/pkg/ci"
	candlesticksCi "github.com/lerenn/cryptellation/svc/candlesticks/pkg/ci"
	exchangesCi "github.com/lerenn/cryptellation/svc/exchanges/pkg/ci"
	indicatorsCi "github.com/lerenn/cryptellation/svc/indicators/pkg/ci"
	ticksCi "github.com/lerenn/cryptellation/svc/ticks/pkg/ci"
	"github.com/spf13/cobra"
)

func generators() map[string]func(context.Context) error {
	return map[string]func(context.Context) error{
		"svc/backtests":    backtestsCi.Generator(client),
		"svc/candlesticks": candlesticksCi.Generator(client),
		"svc/exchanges":    exchangesCi.Generator(client),
		"svc/indicators":   indicatorsCi.Generator(client),
		"svc/ticks":        ticksCi.Generator(client),
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
