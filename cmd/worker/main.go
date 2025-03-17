package main

import (
	"context"
	"os"

	"github.com/lerenn/cryptellation/v1/pkg/telemetry"
	"github.com/lerenn/cryptellation/v1/pkg/telemetry/console"
	"github.com/lerenn/cryptellation/v1/pkg/telemetry/otel"
	"github.com/lerenn/cryptellation/v1/pkg/version"
	"github.com/spf13/cobra"
)

// rootCmd is the worker root command.
var rootCmd = &cobra.Command{
	Use:     "worker",
	Version: version.FullVersion(),
	Short:   "worker - a worker executing cryptellation temporal workflows",
}

func main() {
	var errCode int

	// Init opentelemetry and set it globally
	console.Fallback(otel.NewTelemeter(context.Background(), "cryptellation-worker"))

	// Set commands
	rootCmd.AddCommand(serveCmd)
	addDatabaseCommands(rootCmd)

	// Execute command
	if err := rootCmd.Execute(); err != nil {
		telemetry.L(context.Background()).Errorf("an error occurred: %s", err.Error())
	}

	// Close telemetry
	telemetry.Close(context.Background())

	// Exit with error code
	os.Exit(errCode)
}
