package main

import (
	"context"
	"fmt"

	"dagger.io/dagger"
	"github.com/lerenn/cryptellation/pkg/ci"
	backtestsCi "github.com/lerenn/cryptellation/svc/backtests/pkg/ci"
	candlesticksCi "github.com/lerenn/cryptellation/svc/candlesticks/pkg/ci"
	exchangesCi "github.com/lerenn/cryptellation/svc/exchanges/pkg/ci"
	indicatorsCi "github.com/lerenn/cryptellation/svc/indicators/pkg/ci"
	ticksCi "github.com/lerenn/cryptellation/svc/ticks/pkg/ci"
	"github.com/spf13/cobra"
)

func runServers(cmd *cobra.Command, args []string) {
	// Dependencies
	broker := ci.Nats(client).AsService()
	withBroker := ci.NatsDependency(broker)

	// Cryptellation that will be run as dependencies
	candlesticks := candlesticksCi.RunnerWithDependencies(client, withBroker)

	// Cryptellation
	backtests := backtestsCi.RunnerWithDependencies(client, withBroker, candlesticks.AsService())
	exchanges := exchangesCi.RunnerWithDependencies(client, withBroker)
	indicators := indicatorsCi.RunnerWithDependencies(client, withBroker, candlesticks.AsService())
	ticks := ticksCi.RunnerWithDependencies(client, withBroker)

	// Set tunnel to NATS
	tunnel, err := client.Host().Tunnel(broker).Start(context.Background())
	if err != nil {
		panic(err)
	}
	defer tunnel.Stop(context.Background())

	// Get NATS service address
	srvAddr, err := tunnel.Endpoint(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println("NATS running on host at this address:", srvAddr)

	// Run services
	ci.ExecuteContainersInParallel(context.Background(), []*dagger.Container{
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
