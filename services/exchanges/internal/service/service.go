package service

import (
	"github.com/digital-feather/cryptellation/services/exchanges/internal/adapters/exchanges"
	"github.com/digital-feather/cryptellation/services/exchanges/internal/adapters/exchanges/binance"
	"github.com/digital-feather/cryptellation/services/exchanges/internal/application"
)

func NewApplication() (*application.Application, error) {
	binanceService, err := binance.New()
	if err != nil {
		return nil, err
	}

	services := map[string]exchanges.Adapter{
		exchanges.Binance.Name: binanceService,
	}

	return application.New(services)
}

func newMockApplication() (*application.Application, error) {
	return application.New(map[string]exchanges.Adapter{
		"mock_exchange": MockExchangeService{},
	})
}
