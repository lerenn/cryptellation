package main

import (
	"encoding/json"
	"fmt"
	"os"

	client "github.com/lerenn/cryptellation/clients/go"

	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/version"

	"github.com/spf13/cobra"
)

var (
	globalClient client.Services
	jsonOutput   bool
)

var CryptellationCmd = &cobra.Command{
	Use:     "cryptellation",
	Version: version.FullVersion(),
	Short:   "cryptellation - a simple CLI to manipulate cryptellation services",
	Long: "cryptellation is a simple CLI to manipulate cryptellation services.\n\n" +
		"One can use cryptellation-candlesticks to manage migrations from the terminal and launch the service.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		// Get configuration
		conf := config.LoadNATS()
		if err := conf.Validate(); err != nil {
			return err
		}

		// Create client
		globalClient, err = client.NewServices(conf)
		return err
	},
}

func init() {
	// Flags
	CryptellationCmd.PersistentFlags().BoolVarP(&jsonOutput, "json", "j", false, "Set output to JSON format")

	// Services
	initBacktests(CryptellationCmd)
	initCandlesticks(CryptellationCmd)
	initExchanges(CryptellationCmd)
	initForwardTests(CryptellationCmd)
	initIndicators(CryptellationCmd)
	initTicks(CryptellationCmd)

	// Others
	initStatuses(CryptellationCmd)
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

func displayJSON(jsonObj any) error {
	output, err := json.Marshal(jsonObj)
	if err != nil {
		return err
	}

	fmt.Println(string(output))
	return nil
}
