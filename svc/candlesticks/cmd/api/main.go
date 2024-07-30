package main

import (
	"context"
	"os"

	"cryptellation/pkg/adapters/telemetry"
	"cryptellation/pkg/adapters/telemetry/console"
	"cryptellation/pkg/adapters/telemetry/otel"
	"cryptellation/pkg/version"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:     "api",
	Version: version.FullVersion(),
	Short:   "api - a simple CLI to manipulate candlesticks service",
}

func init() {
	RootCmd.AddCommand(serveCmd)
}

func main() {
	// Init opentelemetry and set it globally
	console.Fallback(otel.NewTelemeter(context.Background(), "cryptellation-candlesticks"))
	defer telemetry.Close(context.TODO())

	// Execute command
	if err := RootCmd.Execute(); err != nil {
		telemetry.L(context.TODO()).Errorf("an error occured: %s", err.Error())
		os.Exit(1)
	}
}
