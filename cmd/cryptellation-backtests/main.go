package main

import (
	"fmt"
	"os"

	"github.com/lerenn/cryptellation/pkg/version"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:     "cryptellation-backtests",
	Version: version.GetFullVersion(),
	Short:   "cryptellation-backtests - a simple CLI to manipulate backtests service",
	Long: "cryptellation-backtests is a simple CLI to manipulate backtests service.\n\n" +
		"One can use cryptellation-backtests to manage migrations from the terminal and launch the service.",
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
