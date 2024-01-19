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
	Use:     "api",
	Version: version.FullVersion(),
	Short:   "api - a simple CLI to manipulate backtests service",
}

func init() {
	RootCmd.AddCommand(serveCmd)
}

func main() {
	// Init opentelemetry and set it globally
	tlr, err := otel.NewTelemeter(context.Background(), "cryptellation-backtests")
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
