package main

import (
	"time"

	cryptellation "github.com/lerenn/cryptellation/clients/go"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/utils"
	backtests "github.com/lerenn/cryptellation/svc/backtests/clients/go"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/account"
)

func main() {
	// Create a cryptellation client
	client, err := cryptellation.NewServices(config.LoadNATS())
	if err != nil {
		panic(err)
	}

	// Create backtest
	b, err := cryptellation.NewBacktest(client, backtests.BacktestCreationPayload{
		Accounts: map[string]account.Account{
			"binance": {
				Balances: map[string]float64{
					"USDT": 1000,
				},
			},
		},
		StartTime: utils.Must(time.Parse(time.RFC3339, "2024-01-01T00:00:00Z")),
		EndTime:   utils.ToReference(utils.Must(time.Parse(time.RFC3339, "2024-01-02T00:00:00Z"))),
	}, &Bot{})
	if err != nil {
		panic(err)
	}

	// Run backtest
	if err := b.Run(); err != nil {
		panic(err)
	}
}
