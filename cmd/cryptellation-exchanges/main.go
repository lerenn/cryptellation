package main

import (
	"fmt"
	"os"

	"github.com/digital-feather/cryptellation/pkg/version"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:     "cryptellation-exchanges",
	Version: version.GetFullVersion(),
	Short:   "cryptellation-exchanges - a simple CLI to manipulate exchanges service",
	Long: "cryptellation-exchanges is a simple CLI to manipulate exchanges service.\n\n" +
		"One can use cryptellation-exchanges to manage migrations from the terminal and launch the service.",
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
