package main

import (
	"context"
	"fmt"
	"os"

	"github.com/lerenn/cryptellation/internal/adapters/telemetry"
	"github.com/lerenn/cryptellation/internal/adapters/telemetry/otel"
	"github.com/lerenn/cryptellation/pkg/version"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:     "cryptellation-candlesticks",
	Version: version.GetFullVersion(),
	Short:   "cryptellation-candlesticks - a simple CLI to manipulate candlesticks service",
	Long: "cryptellation-candlesticks is a simple CLI to manipulate candlesticks service.\n\n" +
		"One can use cryptellation-candlesticks to manage migrations from the terminal and launch the service.",
}

func init() {
	RootCmd.AddCommand(serveCmd)

	addCommandsToMigrationsCmd()
	RootCmd.AddCommand(migrationsCmd)
}

func main() {
	// Init opentelemetry and set it globally
	tlr, err := otel.NewTelemeter(context.Background(), "cryptellation-candlesticks")
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occured when setting telemetry: %s", err.Error())
	}
	defer tlr.Close(context.TODO())

	// Set telemetry globally
	telemetry.Set(tlr)

	// Execute command
	if err := RootCmd.Execute(); err != nil {
		tlr.Logger(context.TODO()).Error(
			fmt.Sprintf("an error occured: %s", err.Error()),
		)
		os.Exit(1)
	}
}
