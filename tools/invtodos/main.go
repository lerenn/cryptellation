package main

import (
	"context"
	"fmt"
	"os"

	"github.com/lerenn/cryptellation/v1/pkg/telemetry"
	"github.com/lerenn/cryptellation/v1/pkg/telemetry/console"
	"github.com/lerenn/cryptellation/v1/pkg/telemetry/otel"
	"github.com/lerenn/cryptellation/v1/pkg/version"
	"github.com/spf13/cobra"
)

var (
	pathFlag string
)

var rootCmd = &cobra.Command{
	Use:          os.Args[0],
	Version:      version.FullVersion(),
	SilenceUsage: true,
	Short:        os.Args[0] + " - a CLI to check invalid todos",
	RunE: func(cmd *cobra.Command, _ []string) error {
		invalidLines, err := checkInvalidTodosOnDir(pathFlag)
		if err != nil {
			return err
		}

		if len(invalidLines) > 0 {
			for _, line := range invalidLines {
				telemetry.L(cmd.Context()).Errorf(line)
			}
			return fmt.Errorf("found %d invalid todos", len(invalidLines))
		}

		return nil
	},
}

func main() {
	var errCode int

	// Init opentelemetry and set it globally
	console.Fallback(otel.NewTelemeter(context.Background(), "cryptellation"))

	// Set flags
	rootCmd.PersistentFlags().StringVarP(&pathFlag, "path", "p", ".", "Set the path to check for invalid todos")

	// Execute command
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		telemetry.L(context.Background()).Errorf("an error occurred: %s", err.Error())
		errCode = 1
	}

	// Close telemetry
	telemetry.Close(context.Background())

	// Exit with error code
	os.Exit(errCode)
}
