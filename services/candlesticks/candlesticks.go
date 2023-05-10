package candlesticks

import (
	db "github.com/lerenn/cryptellation/services/candlesticks/io/db"
	"github.com/lerenn/cryptellation/services/candlesticks/io/exchanges"
)

type candlesticks struct {
	db        db.Port
	exchanges exchanges.Port
}

func New(db db.Port, exchanges exchanges.Port) Interface {
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
