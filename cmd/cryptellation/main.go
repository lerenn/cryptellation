package main

import (
	"fmt"
	"os"

	"github.com/lerenn/cryptellation/clients/go/nats"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/version"
	"github.com/spf13/cobra"
)

var (
	services nats.Services
)

var CryptellationCmd = &cobra.Command{
	Use:     "cryptellation",
	Version: version.FullVersion(),
	Short:   "cryptellation - a simple CLI to manipulate cryptellation services",
	Long: "cryptellation is a simple CLI to manipulate cryptellation services.\n\n" +
		"One can use cryptellation-candlesticks to manage migrations from the terminal and launch the service.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		// Get clients
		services, err = nats.NewServices(config.LoadNATS())
		return err
	},
}

func init() {
	initCandlesticks(CryptellationCmd)
	initExchanges(CryptellationCmd)
	initTicks(CryptellationCmd)
	initInfo(CryptellationCmd)
}

func main() {
	if err := CryptellationCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

// this is compensating this issue: https://github.com/spf13/cobra/issues/252
func executeParentPersistentPreRuns(cmd *cobra.Command, args []string) error {
	root := cmd
	for root.HasParent() {
		root = root.Parent()
	}
	return root.PersistentPreRunE(cmd, args)
}
