package binance

import (
	"context"
	"fmt"
	"time"

	"github.com/lerenn/cryptellation/pkg/adapters/exchanges/binance"
	"github.com/lerenn/cryptellation/svc/exchanges/pkg/exchange"
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

func (ps *Service) Infos(ctx context.Context) (exchange.Exchange, error) {
	exchangeInfos, err := ps.Client.NewExchangeInfoService().Do(ctx)
	if err != nil {
		return exchange.Exchange{}, err
	}

	pairs := make([]string, len(exchangeInfos.Symbols))
	for i, bs := range exchangeInfos.Symbols {
		pairs[i] = fmt.Sprintf("%s-%s", bs.BaseAsset, bs.QuoteAsset)
	}

	exch := binance.Infos
	exch.Pairs = pairs
	exch.LastSyncTime = time.Now()

	return exch, nil
}
