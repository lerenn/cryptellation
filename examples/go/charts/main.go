package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	cryptellation "cryptellation/client"

	"cryptellation/internal/config"
	"cryptellation/pkg/charts"
	"cryptellation/pkg/utils"

	candlesticks "cryptellation/svc/candlesticks/clients/go"
	"cryptellation/svc/candlesticks/pkg/candlestick"
	"cryptellation/svc/candlesticks/pkg/period"

	indicators "cryptellation/svc/indicators/clients/go"
)

func main() {
	http.HandleFunc("/", httpserver)
	_ = http.ListenAndServe(":8081", nil)
}

func httpserver(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request")

	// Create a cryptellation client
	client, err := cryptellation.NewServices(config.LoadNATS())
	if err != nil {
		panic(err)
	}

	// Create a chts generator
	chts := charts.NewGenerator(client)

	// Generate candlesticks
	fmt.Println("Generating candlesticks")
	start := utils.Must(time.Parse(time.RFC3339, "2021-01-01T00:00:00Z"))
	end := utils.Must(time.Parse(time.RFC3339, "2021-01-01T01:00:00Z"))
	c, err := chts.Candlesticks(context.Background(), charts.CandlesticksPayload{
		Name: "BTC-USDT",
		Candlesticks: candlesticks.ReadCandlesticksPayload{
			Exchange: "binance",
			Pair:     "BTC-USDT",
			Period:   period.M1,
			Start:    &start,
			End:      &end,
		},
	})
	if err != nil {
		panic(err)
	}

	// Generate SMAs
	fmt.Println("Generating SMAs")
	payload := charts.SMAPayload{
		SMA: indicators.SMAPayload{
			Exchange:  "binance",
			Pair:      "BTC-USDT",
			Period:    period.M1,
			Start:     start,
			End:       end,
			PriceType: candlestick.PriceTypeIsClose,
		},
	}

	payload.SMA.PeriodNumber = 7
	payload.Color = "red"
	sma7, err := chts.SMA(context.Background(), payload)
	if err != nil {
		panic(err)
	}
	c.Overlap(sma7)

	payload.SMA.PeriodNumber = 25
	payload.Color = "yellow"
	sma25, err := chts.SMA(context.Background(), payload)
	if err != nil {
		panic(err)
	}
	c.Overlap(sma25)

	payload.SMA.PeriodNumber = 99
	payload.Color = "green"
	sma99, err := chts.SMA(context.Background(), payload)
	if err != nil {
		panic(err)
	}
	c.Overlap(sma99)

	if err := c.Render(w); err != nil {
		panic(err)
	}
}
