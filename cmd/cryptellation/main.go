package main

import (
	"context"
	"os"

	cryptellationclient "github.com/lerenn/cryptellation/v1/pkg/client"
	"github.com/lerenn/cryptellation/v1/pkg/telemetry"
	"github.com/lerenn/cryptellation/v1/pkg/telemetry/console"
	"github.com/lerenn/cryptellation/v1/pkg/telemetry/otel"
	"github.com/lerenn/cryptellation/v1/pkg/version"
	"github.com/spf13/cobra"
)

var (
	cryptellationClient cryptellationclient.Client
)

// RootCmd is the CLI root command.
var RootCmd = &cobra.Command{
	Use:     "cryptellation",
	Version: version.FullVersion(),
	Short:   "cryptellation - a CLI to execute cryptellation temporal workflows",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		// Create cryptellation client
		cryptellationClient, err = cryptellationclient.New()
		return err
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		cryptellationClient.Close(cmd.Context())
	},
}

func main() {
	var errCode int

	// Init opentelemetry and set it globally
	console.Fallback(otel.NewTelemeter(context.Background(), "cryptellation"))

	// Set commands
	RootCmd.AddCommand(infoCmd)

	// Execute command
	if err := RootCmd.Execute(); err != nil {
		telemetry.L(context.Background()).Errorf("an error occurred: %s", err.Error())
	}

	// Close telemetry
	telemetry.Close(context.Background())

	// Exit with error code
	os.Exit(errCode)
}
