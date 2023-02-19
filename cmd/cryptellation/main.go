package main

import (
	"context"
	"fmt"
	"os"
	"time"

	candlesticks "github.com/digital-feather/cryptellation/internal/candlesticks/ctrl/nats"
	"github.com/digital-feather/cryptellation/pkg/candlestick"
	"github.com/digital-feather/cryptellation/pkg/config"
	"github.com/digital-feather/cryptellation/pkg/period"
	"github.com/digital-feather/cryptellation/pkg/utils"
)

func run() int {
	client, err := candlesticks.New(config.LoadNATSConfigFromEnv())
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occured when %+v\n", fmt.Errorf("initializing client: %w", err))
		return 255
	}

	list, err := client.ReadCandlesticks(context.Background(), candlesticks.ReadCandlesticksPayload{
		ExchangeName: "binance",
		PairSymbol:   "ETH-USDT",
		Period:       period.H1,
		Start:        utils.ToReference(time.Now().AddDate(0, 0, -7)),
		End:          utils.ToReference(time.Now()),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occured when %+v\n", fmt.Errorf("reading candlesticks: %w", err))
		return 255
	}

	list.Loop(func(t time.Time, cs candlestick.Candlestick) (bool, error) {
		fmt.Println(cs.String())
		return false, nil
	})

	return 0
}

func main() {
	os.Exit(run())
}
