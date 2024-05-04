package main

import (
	"context"
	"os"

	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
	"github.com/lerenn/cryptellation/pkg/adapters/telemetry/console"
	"github.com/lerenn/cryptellation/pkg/adapters/telemetry/otel"
	"github.com/lerenn/cryptellation/pkg/version"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:     "data",
	Version: version.FullVersion(),
	Short:   "data - a simple CLI to manipulate backtests service",
}

func init() {
	addCommandsToMigrationsCmd()
	RootCmd.AddCommand(migrationsCmd)
}

func main() {
	// Init opentelemetry and set it globally
	console.Fallback(otel.NewTelemeter(context.Background(), "cryptellation-backtests"))
	defer telemetry.Close(context.TODO())

	// Execute command
	if err := RootCmd.Execute(); err != nil {
		telemetry.L(context.TODO()).Errorf("an error occured: %s", err.Error())
		os.Exit(1)
	}
}
