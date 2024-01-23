package main

import (
	"context"
	"fmt"
	"os"
	"strings"

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
	// Set broker
	broker := ci.Nats(client).AsService()
	withBroker := ci.NatsDependency(broker)
	stop := writeBrokerAccess(broker)
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

func writeBrokerAccess(broker *dagger.Service) func(ctx context.Context) (*dagger.Service, error) {
	// Set tunnel to NATS
	tunnel, err := client.Host().Tunnel(broker).Start(context.Background())
	if err != nil {
		panic(err)
	}

	// Get NATS service address
	srvAddr, err := tunnel.Endpoint(context.Background())
	if err != nil {
		panic(err)
	}

	// Write server address in a file
	f, err := os.Create(".env")
	if err != nil {
		panic(err)
	}

	// Split address in host/port
	srvAddrSplit := strings.Split(srvAddr, ":")

	str := fmt.Sprintf("NATS_HOST=%s\n", srvAddrSplit[0])
	str += fmt.Sprintf("NATS_PORT=%s\n", srvAddrSplit[1])
	_, _ = f.WriteString(str)

	return tunnel.Stop
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
