package main

import (
	"context"

	"github.com/lerenn/cryptellation/pkg/ci"
	"github.com/lerenn/cryptellation/pkg/ci/publish"
	backtestsCi "github.com/lerenn/cryptellation/svc/backtests/pkg/ci"
	candlesticksCi "github.com/lerenn/cryptellation/svc/candlesticks/pkg/ci"
	exchangesCi "github.com/lerenn/cryptellation/svc/exchanges/pkg/ci"
	indicatorsCi "github.com/lerenn/cryptellation/svc/indicators/pkg/ci"
	ticksCi "github.com/lerenn/cryptellation/svc/ticks/pkg/ci"
	"github.com/spf13/cobra"
)

var (
	tagsFlag []string
)

func dockerImagePublishers() map[string]func(ctx context.Context) error {
	return map[string]func(ctx context.Context) error{
		pathCmdCli: publish.PublishDockerImage(
			backtestsCi.Runner(client),
			pathCmdCli[1:],
			"lerenn/cryptellation-cli"),
		pathCmdTui: publish.PublishDockerImage(
			backtestsCi.Runner(client),
			pathCmdTui[1:],
			"lerenn/cryptellation-tui"),
		pathSvcBacktests: publish.PublishDockerImage(
			backtestsCi.Runner(client),
			pathSvcBacktests[1:],
			"lerenn/cryptellation-backtests"),
		pathSvcCandlesticks: publish.PublishDockerImage(
			candlesticksCi.Runner(client),
			pathSvcCandlesticks[1:],
			"lerenn/cryptellation-candlesticks"),
		pathSvcExchanges: publish.PublishDockerImage(
			exchangesCi.Runner(client),
			pathSvcExchanges[1:],
			"lerenn/cryptellation-exchanges"),
		pathSvcIndicators: publish.PublishDockerImage(
			indicatorsCi.Runner(client),
			pathSvcIndicators[1:],
			"lerenn/cryptellation-indicators"),
		pathSvcTicks: publish.PublishDockerImage(
			ticksCi.Runner(client),
			pathSvcTicks[1:],
			"lerenn/cryptellation-ticks"),
	}
}

func runPublishers(cmd *cobra.Command, args []string) error {
	modules := removeLeadingSlash(pathModules)
	if err := publish.GitTagAndPush(modules, tagsFlag); err != nil {
		return err
	}

	ci.ExecuteInParallel(context.Background(), filterWithPath(dockerImagePublishers())...)
	return nil
}

var publishCmd = &cobra.Command{
	Use:     "publish",
	Aliases: []string{"p"},
	Short:   "Execute publish step of the CI",
	RunE:    runPublishers,
}

func addPublishCmdTo(cmd *cobra.Command) {
	publishCmd.Flags().StringArrayVarP(
		&tagsFlag, "tags", "t", []string{},
		"Tags used to tag services (with <path|bump> where "+
			"'path' is the module path, 'empty' for codebase or '*' for all modules"+
			"and 'bump' is 'major', 'minor' or 'fix')")
	cmd.AddCommand(publishCmd)
}
