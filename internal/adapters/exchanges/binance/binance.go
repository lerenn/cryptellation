package binance

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	client "github.com/adshao/go-binance/v2"
	"github.com/lerenn/cryptellation/internal/adapters/exchanges/binance/entities"
	"github.com/lerenn/cryptellation/internal/components/candlesticks/ports/exchanges"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/pkg/models/exchange"
	"github.com/lerenn/cryptellation/pkg/models/pair"
	"github.com/lerenn/cryptellation/pkg/models/tick"
)

var (
	Infos = exchange.Exchange{
		Name: entities.Name,
		PeriodsSymbols: []string{
			"M1", "M3", "M5", "M15", "M30",
			"H1", "H2", "H4", "H6", "H8", "H12",
			"D1", "D3",
			"W1",
		},
		Fees: 0.1,
	}
)

type Service struct {
	client *client.Client
}

func New() (*Service, error) {
	// Get config
	config := config.LoadBinance()

	// Return service
	return &Service{
		client: client.NewClient(
			config.ApiKey,
			config.SecretKey),
	}, config.Validate()
}

func (ps *Service) Infos(ctx context.Context) (exchange.Exchange, error) {
	exchangeInfos, err := ps.client.NewExchangeInfoService().Do(ctx)
	if err != nil {
		return exchange.Exchange{}, err
	}

	pairSymbols := make([]string, len(exchangeInfos.Symbols))
	for i, bs := range exchangeInfos.Symbols {
		pairSymbols[i] = fmt.Sprintf("%s-%s", bs.BaseAsset, bs.QuoteAsset)
	}

	exch := Infos
	exch.PairsSymbols = pairSymbols
	exch.LastSyncTime = time.Now()

	return exch, nil
}

func (s *Service) ListenSymbol(symbol string) (chan tick.Tick, chan struct{}, error) {
	binanceSymbol, err := toBinanceSymbol(symbol)
	if err != nil {
		return nil, nil, err
	}

	tickChan := make(chan tick.Tick)
	_, stop, err := client.WsBookTickerServe(binanceSymbol, func(event *client.WsBookTickerEvent) {
		ask, err := strconv.ParseFloat(event.BestAskPrice, 64)
		if err != nil {
			log.Println(err)
			return
		}

		bid, err := strconv.ParseFloat(event.BestBidPrice, 64)
		if err != nil {
			log.Println(err)
			return
		}

		t := tick.Tick{
			Time:       time.Now().UTC(),
			Exchange:   "binance",
			PairSymbol: symbol,
			Price:      float64(ask+bid) / 2,
		}

		// Send it to tick channel
		select {
		case tickChan <- t:
		default:
			log.Printf("Dropped %q tick from binance adapter\n", symbol)
		}

	}, nil)

	// TODO: manage when error or done

	return tickChan, stop, err
}

func toBinanceSymbol(symbol string) (string, error) {
	base, quote, err := pair.ParsePairSymbol(symbol)
	return fmt.Sprintf("%s%s", base, quote), err
}

func (s *Service) GetCandlesticks(ctx context.Context, payload exchanges.GetCandlesticksPayload) (*candlestick.List, error) {
	s.client.Debug = true

	service := s.client.NewKlinesService()

	// Set symbol
	service.Symbol(entities.BinanceSymbol(payload.PairSymbol))

	// Set interval
	binanceInterval, err := entities.PeriodToInterval(payload.Period)
	if err != nil {
		return nil, entities.WrapError(err)
	}
	service.Interval(binanceInterval)

	// Set start and end time
	service.StartTime(entities.TimeToKLineTime(payload.Start))
	service.EndTime(entities.TimeToKLineTime(payload.End))

	// Set limit
	if payload.Limit > 0 {
		service.Limit(payload.Limit)
	}

	// Get KLines
	kl, err := service.Do(ctx)
	if err != nil {
		return nil, entities.WrapError(err)
	}

	// Change them to right format
	return entities.KLinesToCandlesticks(payload.PairSymbol, payload.Period, kl, time.Now())
}
