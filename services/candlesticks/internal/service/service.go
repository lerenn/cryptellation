package service

import (
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/adapters/exchanges"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/adapters/exchanges/binance"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/application"
)

func NewApplication() (*application.Application, error) {
	binanceService, err := binance.New()
	if err != nil {
		return nil, err
	}

	services := map[string]exchanges.Adapter{
		binance.Name: binanceService,
	}

	return application.New(services)
}

func newMockApplication() (*application.Application, error) {
	services := map[string]exchanges.Adapter{
		"mock_exchange": &MockExchangeService{},
	}

	return application.New(services)
}
