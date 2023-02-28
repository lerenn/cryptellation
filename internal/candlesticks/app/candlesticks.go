package app

import (
	db "github.com/digital-feather/cryptellation/internal/candlesticks/app/ports/db"
	"github.com/digital-feather/cryptellation/internal/candlesticks/app/ports/exchanges"
)

type candlesticks struct {
	db        db.Port
	exchanges exchanges.Port
}

func New(db db.Port, exchanges exchanges.Port) Controller {
	if db == nil {
		panic("nil db")
	}

	if exchanges == nil {
		panic("nil exchanges")
	}

	return candlesticks{
		db:        db,
		exchanges: exchanges,
	}
}
