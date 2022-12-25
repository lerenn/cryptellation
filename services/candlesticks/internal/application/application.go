package application

import (
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/application/operators/candlesticks"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/application/ports/db"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/application/ports/exchanges"
)

type Application struct {
	Candlesticks candlesticks.Operator
}

func New(db db.Adapter, services map[string]exchanges.Adapter) (*Application, error) {
	return &Application{
		Candlesticks: candlesticks.New(db, services),
	}, nil
}
