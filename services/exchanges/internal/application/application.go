package application

import (
	"github.com/digital-feather/cryptellation/services/exchanges/internal/application/operators/exchanges"
	"github.com/digital-feather/cryptellation/services/exchanges/internal/application/ports/db"
	exchangesAdapters "github.com/digital-feather/cryptellation/services/exchanges/internal/application/ports/exchanges"
)

type Application struct {
	Exchanges exchanges.Operator
}

func New(db db.Adapter, services map[string]exchangesAdapters.Adapter) (*Application, error) {
	return &Application{
		Exchanges: exchanges.New(db, services),
	}, nil
}
