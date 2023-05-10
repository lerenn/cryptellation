package main

import (
	"fmt"
	"os"

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
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "An error occured: %s", err.Error())
		os.Exit(1)
	}
}
