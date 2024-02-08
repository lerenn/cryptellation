package main

import (
	"context"

	"dagger.io/dagger"
	backtestsCi "github.com/lerenn/cryptellation/svc/backtests/pkg/ci"
	candlesticksCi "github.com/lerenn/cryptellation/svc/candlesticks/pkg/ci"
	exchangesCi "github.com/lerenn/cryptellation/svc/exchanges/pkg/ci"
	indicatorsCi "github.com/lerenn/cryptellation/svc/indicators/pkg/ci"
	ticksCi "github.com/lerenn/cryptellation/svc/ticks/pkg/ci"
	"github.com/lerenn/cryptellation/tools/pkg/ci"
	"github.com/spf13/cobra"
)

func runServers(cmd *cobra.Command, args []string) {
	// Set broker
	broker := ci.Nats(client).AsService()
	withBroker := ci.NatsDependency(broker)
	stopBrokerTunnel := ci.ExposeOnLocalPort(client, broker, dagger.PortForward{
		Frontend: 4222,
		Backend:  4222,
	})
	defer stopBrokerTunnel(context.Background()) //nolint: errcheck, no need to check error here

	uptrace, otelcollector := ci.Uptrace(client)
	withOtelCollector := ci.OtelCollectorDependency(otelcollector)
	stopUptraceTunnel := ci.ExposeOnLocalPort(client, uptrace, dagger.PortForward{
		Frontend: 4318,
		Backend:  4318,
	})
	defer stopUptraceTunnel(context.TODO())

	// Cryptellation
	candlesticks := candlesticksCi.RunnerWithDependencies(client, withBroker, withOtelCollector)
	backtests := backtestsCi.RunnerWithDependencies(client, withBroker, withOtelCollector)
	exchanges := exchangesCi.RunnerWithDependencies(client, withBroker, withOtelCollector)
	indicators := indicatorsCi.RunnerWithDependencies(client, withBroker, withOtelCollector)
	ticks := ticksCi.RunnerWithDependencies(client, withBroker, withOtelCollector)

	// Run services
	ci.ExecuteContainersInParallel(context.Background(), []*dagger.Container{
		candlesticks,
		backtests,
		exchanges,
		indicators,
		ticks,
	})
}

var serveCmd = &cobra.Command{
	Use:     "serve",
	Aliases: []string{"s"},
	Short:   "Execute serve step of the CI",
	Run:     runServers,
}

func addServeCmdTo(cmd *cobra.Command) {
	cmd.AddCommand(serveCmd)
}
