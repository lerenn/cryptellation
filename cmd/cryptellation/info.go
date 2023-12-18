package main

import (
	"context"
	"fmt"

	"gopkg.in/yaml.v2"

	client "github.com/lerenn/cryptellation/clients/go"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"i"},
	Short:   "Read info from all cryptellation services",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		servicesInfo := make(map[string]client.ServiceInfo)

		// List all functions to get info
		servicesInfoFuncs := map[string]func(context.Context) (client.ServiceInfo, error){
			"backtests":    services.Backtests.ServiceInfo,
			"candlesticks": services.Candlesticks.ServiceInfo,
			"exchanges":    services.Exchanges.ServiceInfo,
			"indicators":   services.Indicators.ServiceInfo,
			"ticks":        services.Ticks.ServiceInfo,
		}

		// Get info from each service
		for name, fn := range servicesInfoFuncs {
			s, err := fn(context.TODO())
			if err != nil {
				return err
			}
			servicesInfo[name] = s
		}

		// Display results
		bString, err := yaml.Marshal(servicesInfo)
		if err != nil {
			return err
		}
		fmt.Println(string(bString))

		return nil
	},
}

func initInfo(rootCmd *cobra.Command) {
	rootCmd.AddCommand(infoCmd)
}
