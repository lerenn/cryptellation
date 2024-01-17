package main

import (
	"context"

	"dagger.io/dagger"
	"github.com/lerenn/cryptellation/pkg/ci"
	backtestsCi "github.com/lerenn/cryptellation/svc/backtests/pkg/ci"
	candlesticksCi "github.com/lerenn/cryptellation/svc/candlesticks/pkg/ci"
	exchangesCi "github.com/lerenn/cryptellation/svc/exchanges/pkg/ci"
	indicatorsCi "github.com/lerenn/cryptellation/svc/indicators/pkg/ci"
	ticksCi "github.com/lerenn/cryptellation/svc/ticks/pkg/ci"
	"github.com/spf13/cobra"
)

func binaryBuilders() map[string]*dagger.Container {
	return map[string]*dagger.Container{
		"svc/backtests":    backtestsCi.BuildBinary(client),
		"svc/candlesticks": candlesticksCi.BuildBinary(client),
		"svc/exchanges":    exchangesCi.BuildBinary(client),
		"svc/indicators":   indicatorsCi.BuildBinary(client),
		"svc/ticks":        ticksCi.BuildBinary(client),
	}
}

func runBuilders(cmd *cobra.Command, args []string) {
	ci.ExecuteContainersInParallel(context.Background(), filterWithPath(binaryBuilders()))
}

var buildCmd = &cobra.Command{
	Use:     "build",
	Aliases: []string{"b"},
	Short:   "Execute build step of the CI",
	Run:     runBuilders,
}

func addBuildCmdTo(cmd *cobra.Command) {
	cmd.AddCommand(buildCmd)
}
