package main

import (
	"fmt"
	"os"

	"github.com/lerenn/cryptellation/pkg/version"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:     "cryptellation-indicators",
	Version: version.GetFullVersion(),
	Short:   "cryptellation-indicators - a simple CLI to manipulate indicators service",
	Long: "cryptellation-indicators is a simple CLI to manipulate indicators service.\n\n" +
		"One can use cryptellation-indicators to manage migrations from the terminal and launch the service.",
}

func init() {
	RootCmd.AddCommand(serveCmd)

	addCommandsToMigrationsCmd()
	RootCmd.AddCommand(migrationsCmd)
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "An error occured: %s", err.Error())
		os.Exit(1)
	}
}
