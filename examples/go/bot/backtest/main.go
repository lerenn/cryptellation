package main

import (
	"context"
	"time"

	cryptellation "github.com/lerenn/cryptellation/clients/go"

	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
	"github.com/lerenn/cryptellation/pkg/adapters/telemetry/console"
	"github.com/lerenn/cryptellation/pkg/adapters/telemetry/otel"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/models/account"
	"github.com/lerenn/cryptellation/pkg/utils"

	backtests "github.com/lerenn/cryptellation/client/clients/go"

	"github.com/lerenn/cryptellation/examples/go/bot"
)

func main() {
	// Init opentelemetry and set it globally
	console.Fallback(otel.NewTelemeter(context.Background(), "cryptellation-backtests"))
	defer telemetry.Close(context.Background())

	// Create a cryptellation client
	client, err := cryptellation.NewServices(config.LoadNATS())
	if err != nil {
		panic(err)
	}

	// Create backtest
	b, err := cryptellation.NewBacktest(
		context.Background(),
		client,
		backtests.BacktestCreationPayload{
			Accounts: map[string]account.Account{
				"binance": {
					Balances: map[string]float64{
						"USDT": 1000,
					},
				},
			},
			StartTime: utils.Must(time.Parse(time.RFC3339, "2024-01-01T00:00:00Z")),
			EndTime:   utils.ToReference(utils.Must(time.Parse(time.RFC3339, "2024-01-02T00:00:00Z"))),
		},
		&bot.Bot{})
	if err != nil {
		panic(err)
	}

	// Run backtest
	if err := b.Run(context.Background()); err != nil {
		panic(err)
	}
}
