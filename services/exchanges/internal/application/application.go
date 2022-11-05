package application

import (
	"github.com/digital-feather/cryptellation/services/exchanges/internal/adapters/db/sql"
	exchangesAdapters "github.com/digital-feather/cryptellation/services/exchanges/internal/adapters/exchanges"
	"github.com/digital-feather/cryptellation/services/exchanges/internal/application/exchanges"
)

type Application struct {
	Exchanges exchanges.Operator
}

func New(services map[string]exchangesAdapters.Adapter) (*Application, error) {
	repository, err := sql.New()
	if err != nil {
		return nil, err
	}

	return &Application{
		Exchanges: exchanges.New(repository, services),
	}, nil
}
