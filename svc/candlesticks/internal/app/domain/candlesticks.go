package domain

import (
	db "cryptellation/svc/candlesticks/internal/app/ports/db"
	"cryptellation/svc/candlesticks/internal/app/ports/exchanges"
)

type Candlesticks struct {
	db        db.Port
	exchanges exchanges.Port
}

func New(db db.Port, exchanges exchanges.Port) Candlesticks {
	if db == nil {
		panic("nil db")
	}

	if exchanges == nil {
		panic("nil exchanges")
	}

	return Candlesticks{
		db:        db,
		exchanges: exchanges,
	}
}
