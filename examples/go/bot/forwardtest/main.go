package main

import (
	"context"

	cryptellation "cryptellation/client"

	"cryptellation/pkg/adapters/telemetry"
	"cryptellation/pkg/adapters/telemetry/console"
	"cryptellation/pkg/adapters/telemetry/otel"
	"cryptellation/pkg/config"
	"cryptellation/pkg/models/account"

	"cryptellation/svc/forwardtests/pkg/forwardtest"

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

	// Create forwardtest
	b, err := cryptellation.NewForwardTest(
		context.Background(),
		client,
		forwardtest.NewPayload{
			Accounts: map[string]account.Account{
				"binance": {
					Balances: map[string]float64{
						"USDT": 1000,
					},
				},
			},
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
