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
	stop := ci.ExposeOnLocalPort(client, broker, dagger.PortForward{
		Frontend: 4222,
		Backend:  4222,
	})
	defer stop(context.Background()) //nolint: errcheck, no need to check error here

	// Cryptellation that will be run as dependencies
	candlesticks := candlesticksCi.RunnerWithDependencies(client, withBroker)

	// Cryptellation
	backtests := backtestsCi.RunnerWithDependencies(client, withBroker, candlesticks.AsService())
	exchanges := exchangesCi.RunnerWithDependencies(client, withBroker)
	indicators := indicatorsCi.RunnerWithDependencies(client, withBroker, candlesticks.AsService())
	ticks := ticksCi.RunnerWithDependencies(client, withBroker)

	// Run services
	ci.ExecuteContainersInParallel(context.Background(), []*dagger.Container{
		backtests,
		exchanges,
		indicators,
		ticks,
	})
}

func runTest(cmd *cobra.Command, args []string) {
	broker := ci.Nats(client).AsService()
	withBroker := ci.NatsDependency(broker)

	uptrace, otelcollector := ci.Uptrace(client)
	stop := ci.ExposeOnLocalPort(client, uptrace, dagger.PortForward{
		Frontend: 4318,
		Backend:  4318,
	})
	defer stop(context.TODO())

	candlesticks := candlesticksCi.RunnerWithDependencies(client, withBroker).
		WithServiceBinding("otelco", otelcollector).
		WithEnvVariable("OPENTELEMETRY_GRPC_ENDPOINT", "otelco:4317")

	ci.ExecuteContainersInParallel(context.Background(), []*dagger.Container{
		candlesticks,
	})
}

var serveCmd = &cobra.Command{
	Use:     "serve",
	Aliases: []string{"s"},
	Short:   "Execute serve step of the CI",
	Run:     runTest,
}

func addServeCmdTo(cmd *cobra.Command) {
	cmd.AddCommand(serveCmd)
}
