package binance

import (
	"context"
	"fmt"
	"strconv"
	"time"

	client "github.com/adshao/go-binance/v2"
	"github.com/lerenn/cryptellation/pkg/adapters/exchanges/binance"
	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/pair"
	"github.com/lerenn/cryptellation/svc/ticks/pkg/tick"
)

type Service struct {
	*binance.Service
}

func New() (*Service, error) {
	s, err := binance.New()
	return &Service{
		Service: s,
	}, err
}

func (s *Service) ListenSymbol(ctx context.Context, symbol string) (chan tick.Tick, chan struct{}, error) {
	binanceSymbol, err := toBinanceSymbol(symbol)
	if err != nil {
		return nil, nil, err
	}

	tickChan := make(chan tick.Tick)
	_, stop, err := client.WsBookTickerServe(binanceSymbol, func(event *client.WsBookTickerEvent) {
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
			telemetry.L(ctx).Infof("Dropped %q tick from binance adapter\n", symbol)
		}

	}, nil)

	// TODO: manage when error or done

	return tickChan, stop, err
}

func toBinanceSymbol(symbol string) (string, error) {
	base, quote, err := pair.ParsePair(symbol)
	return fmt.Sprintf("%s%s", base, quote), err
}
