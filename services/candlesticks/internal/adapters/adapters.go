package adapters

import (
	dbAdapters "github.com/digital-feather/cryptellation/services/candlesticks/internal/adapters/db"
	exchangesAdapter "github.com/digital-feather/cryptellation/services/candlesticks/internal/adapters/exchanges"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/components/candlesticks/db"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/components/candlesticks/exchanges"
)

type Adapters struct {
	Exchanges exchanges.Port
	Database  db.Port
}

func New(c Config) (Adapters, error) {
	exchanges, err := exchangesAdapter.New(c.Exchanges)
	if err != nil {
		return Adapters{}, err
	}

	db, err := dbAdapters.New(c.Database)
	if err != nil {
		return Adapters{}, err
	}

	return Adapters{
		Exchanges: exchanges,
		Database:  db,
	}, nil
}
