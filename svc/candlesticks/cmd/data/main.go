package main

import (
	"context"
	"fmt"
	"os"

	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
	"github.com/lerenn/cryptellation/pkg/adapters/telemetry/otel"
	"github.com/lerenn/cryptellation/pkg/version"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:     "cryptellation-candlesticks",
	Version: version.FullVersion(),
	Short:   "cryptellation-candlesticks - a simple CLI to manipulate candlesticks service data.",
}

func init() {
	addCommandsToMigrationsCmd()
	RootCmd.AddCommand(migrationsCmd)
}

func main() {
	// Init opentelemetry and set it globally
	tlr, err := otel.NewTelemeter(context.Background(), "cryptellation-candlesticks")
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occured when setting telemetry: %s\n", err.Error())
	} else {
		// Close when exiting
		defer tlr.Close(context.TODO())

		// Set telemetry globally
		telemetry.Set(tlr)
	}

	// Execute command
	if err := RootCmd.Execute(); err != nil {
		tlr.Logger(context.TODO()).Error(
			fmt.Sprintf("an error occured: %s", err.Error()),
		)
		os.Exit(1)
	}
}
