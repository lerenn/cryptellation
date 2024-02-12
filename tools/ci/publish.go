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

func publishers() map[string]func(ctx context.Context) error {
	return map[string]func(ctx context.Context) error{
		"svc/backtests": ci.PublishDockerImage(
			backtestsCi.Runner(client),
			"svc/backtests",
			"lerenn/cryptellation-backtests"),
		"svc/candlesticks": ci.PublishDockerImage(
			candlesticksCi.Runner(client),
			"svc/candlesticks",
			"lerenn/cryptellation-candlesticks"),
		"svc/exchanges": ci.PublishDockerImage(
			exchangesCi.Runner(client),
			"svc/exchanges",
			"lerenn/cryptellation-exchanges"),
		"svc/indicators": ci.PublishDockerImage(
			indicatorsCi.Runner(client),
			"svc/indicators",
			"lerenn/cryptellation-indicators"),
		"svc/ticks": ci.PublishDockerImage(
			ticksCi.Runner(client),
			"svc/ticks",
			"lerenn/cryptellation-ticks"),
	}
}

func runPublishers(cmd *cobra.Command, args []string) {
	ci.ExecuteInParallel(context.Background(), filterWithPath(publishers())...)
}

var publishCmd = &cobra.Command{
	Use:     "publish",
	Aliases: []string{"p"},
	Short:   "Execute publish step of the CI",
	Run:     runPublishers,
}

func addPublishCmdTo(cmd *cobra.Command) {
	cmd.AddCommand(publishCmd)
}
