package binance

import (
	"context"
	"fmt"
	"time"

	client "github.com/adshao/go-binance/v2"
	"github.com/digital-feather/cryptellation/pkg/config"
	"github.com/digital-feather/cryptellation/pkg/exchange"
)

var (
	Infos = exchange.Exchange{
		Name: "binance",
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

func New(c config.Binance) (*Service, error) {
	return &Service{
		client: client.NewClient(
			c.ApiKey,
			c.SecretKey),
	}, c.Validate()
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
