package main

import (
	"context"

	cryptellation "cryptellation/client"

	"cryptellation/internal/adapters/telemetry"
	"cryptellation/internal/adapters/telemetry/console"
	"cryptellation/internal/adapters/telemetry/otel"
	"cryptellation/internal/config"
	"cryptellation/pkg/models/account"

	"cryptellation/svc/forwardtests/pkg/forwardtest"

	"cryptellation/examples/go/bot"
)

func main() {
	// Init opentelemetry and set it globally
	console.Fallback(otel.NewTelemeter(context.Background(), "cryptellation-backtests"))
	defer telemetry.Close(context.TODO())

	// Create a cryptellation client
	client, err := cryptellation.NewServices(config.LoadNATS())
	if err != nil {
		panic(err)
	}

	// Create forwardtest
	b, err := cryptellation.NewForwardTest(client, forwardtest.NewPayload{
		Accounts: map[string]account.Account{
			"binance": {
				Balances: map[string]float64{
					"USDT": 1000,
				},
			},
		},
	}, &bot.Bot{})
	if err != nil {
		panic(err)
	}

	// Run backtest
	if err := b.Run(); err != nil {
		panic(err)
	}
}
