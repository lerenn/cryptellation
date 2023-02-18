package app

import (
	db "github.com/digital-feather/cryptellation/internal/candlesticks/app/ports/db"
	"github.com/digital-feather/cryptellation/internal/candlesticks/app/ports/exchanges"
)

// Test interface implementation
var _ Port = (*Component)(nil)

type Component struct {
	db        db.Port
	exchanges exchanges.Port
}

func New(db db.Port, exchanges exchanges.Port) Component {
	if db == nil {
		panic("nil db")
	}

	if exchanges == nil {
		panic("nil exchanges")
	}

	return Component{
		db:        db,
		exchanges: exchanges,
	}
}
