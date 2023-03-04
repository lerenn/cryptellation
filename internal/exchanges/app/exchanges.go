package app

import (
	"github.com/digital-feather/cryptellation/internal/exchanges/app/ports/db"
	"github.com/digital-feather/cryptellation/internal/exchanges/app/ports/exchanges"
)

type Exchanges struct {
	db        db.Port
	exchanges exchanges.Port
}

func New(db db.Port, exchanges exchanges.Port) Exchanges {
	if db == nil {
		panic("nil db")
	}

	if exchanges == nil {
		panic("nil exchanges")
	}

	return Exchanges{
		db:        db,
		exchanges: exchanges,
	}
}
