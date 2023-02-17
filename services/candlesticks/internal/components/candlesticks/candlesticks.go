package candlesticks

import (
	db "github.com/digital-feather/cryptellation/services/candlesticks/internal/components/candlesticks/db"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/components/candlesticks/exchanges"
)

// Test interface implementation
var _ Port = (*Candlesticks)(nil)

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
