package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/lerenn/cryptellation/v1/clients/go/client"
	"github.com/lerenn/cryptellation/v1/pkg/telemetry"
	"github.com/lerenn/cryptellation/v1/pkg/telemetry/console"
	"github.com/lerenn/cryptellation/v1/pkg/telemetry/otel"
	"github.com/lerenn/cryptellation/v1/pkg/version"
	"github.com/spf13/cobra"
)

var (
	jsonOutput          bool
	cryptellationClient client.Client
)

// rootCmd is the CLI root command.
var rootCmd = &cobra.Command{
	Use:     os.Args[0],
	Version: version.FullVersion(),
	Short:   os.Args[0] + " - a CLI to execute cryptellation temporal workflows",
	PersistentPreRunE: func(_ *cobra.Command, _ []string) (err error) {
		// Create cryptellation client
		cryptellationClient, err = client.NewClient()
		return err
	},
	PersistentPostRun: func(cmd *cobra.Command, _ []string) {
		cryptellationClient.Close(cmd.Context())
	},
}

func main() {
	var errCode int

	// Init opentelemetry and set it globally
	console.Fallback(otel.NewTelemeter(context.Background(), "cryptellation"))

	// Set flags
	rootCmd.PersistentFlags().BoolVarP(&jsonOutput, "json", "j", false, "Set output to JSON format")

	// Set commands
	addCandlesticksCommands()
	addExchangesCommands()
	rootCmd.AddCommand(infoCmd)
	addTicksCommands()

	// Execute command
	if err := rootCmd.Execute(); err != nil {
		telemetry.L(context.Background()).Errorf("an error occurred: %s", err.Error())
	}

	// Close telemetry
	telemetry.Close(context.Background())

	// Exit with error code
	os.Exit(errCode)
}

func displayJSON(ctx context.Context, jsonObj any) error {
	output, err := json.Marshal(jsonObj)
	if err != nil {
		return err
	}

	telemetry.L(ctx).Info(string(output))
	return nil
}
