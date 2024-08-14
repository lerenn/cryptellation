package binance

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"cryptellation/internal/adapters/exchanges/binance"
	"cryptellation/internal/adapters/telemetry"
	"cryptellation/internal/config"

	"cryptellation/svc/candlesticks/pkg/pair"

	"cryptellation/svc/ticks/pkg/tick"

	client "github.com/adshao/go-binance/v2"
)

type Service struct {
	*binance.Service
}

func New() (*Service, error) {
	s, err := binance.New(config.LoadBinanceTest())
	return &Service{
		Service: s,
	}, err
}

func (s *Service) ListenSymbol(ctx context.Context, symbol string) (chan tick.Tick, chan struct{}, error) {
	binanceSymbol, err := toBinanceSymbol(symbol)
	if err != nil {
		return nil, nil, err
	}

	var lastBid, lastAsk string
	tickChan := make(chan tick.Tick, 64)
	_, stop, err := client.WsBookTickerServe(binanceSymbol, func(event *client.WsBookTickerEvent) {
		// Skip if same price as last tick
		if event.BestAskPrice == lastAsk && event.BestBidPrice == lastBid {
			return
		}
		lastAsk = event.BestAskPrice
		lastBid = event.BestBidPrice

		ask, err := strconv.ParseFloat(event.BestAskPrice, 64)
		if err != nil {
			telemetry.L(ctx).Error(err.Error())
			return
		}

		bid, err := strconv.ParseFloat(event.BestBidPrice, 64)
		if err != nil {
			telemetry.L(ctx).Info(err.Error())
			return
		}

		t := tick.Tick{
			Time:     time.Now().UTC(),
			Exchange: "binance",
			Pair:     symbol,
			Price:    float64(ask+bid) / 2,
		}

		// Send it to tick channel
		select {
		case tickChan <- t:
		default:
			telemetry.L(ctx).Warningf("Dropped %q tick from binance adapter", symbol)
		}

	}, nil)

	// TODO: manage when error or done

	return tickChan, stop, err
}

func toBinanceSymbol(symbol string) (string, error) {
	base, quote, err := pair.ParsePair(symbol)
	return fmt.Sprintf("%s%s", base, quote), err
}
