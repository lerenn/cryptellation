package application

import (
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/adapters/db/sql"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/adapters/exchanges"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/application/candlesticks"
)

type Application struct {
	Candlesticks candlesticks.Operator
}

func New(services map[string]exchanges.Port) (*Application, error) {
	repository, err := sql.New()
	if err != nil {
		return nil, err
	}

	return &Application{
		Candlesticks: candlesticks.New(repository, services),
	}, nil
}
