package domain

import (
	"cryptellation/svc/exchanges/internal/app/ports/db"
	"cryptellation/svc/exchanges/internal/app/ports/exchanges"
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
