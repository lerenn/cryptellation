package service

import (
	"github.com/digital-feather/cryptellation/services/ticks/internal/adapters/exchanges"
	"github.com/digital-feather/cryptellation/services/ticks/internal/adapters/exchanges/binance"
	"github.com/digital-feather/cryptellation/services/ticks/internal/application"
)

func NewApplication() (*application.Application, error) {
	exchanges, err := instanciateExchanges()
	if err != nil {
		return nil, err
	}

	return application.New(exchanges)
}

func instanciateExchanges() (map[string]exchanges.Adapter, error) {
	var err error
	exchanges := make(map[string]exchanges.Adapter)

	exchanges[binance.Name], err = binance.New()
	if err != nil {
		return nil, err
	}

	return exchanges, nil
}

func NewMockedApplication() (*application.Application, error) {
	services := map[string]exchanges.Adapter{
		MockExchangeName: &MockExchangeService{},
	}

	return application.New(services)
}
