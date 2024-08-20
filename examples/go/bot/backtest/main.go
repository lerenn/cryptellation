package main

import (
	"context"
	"time"

	cryptellation "cryptellation/client"

	"cryptellation/internal/config"
	"cryptellation/pkg/adapters/telemetry"
	"cryptellation/pkg/adapters/telemetry/console"
	"cryptellation/pkg/adapters/telemetry/otel"
	"cryptellation/pkg/models/account"
	"cryptellation/pkg/utils"

	backtests "cryptellation/svc/backtests/clients/go"

	"cryptellation/examples/go/bot"
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
