package candlesticks

import (
	db "github.com/digital-feather/cryptellation/internal/components/candlesticks/ports/db"
	"github.com/digital-feather/cryptellation/internal/components/candlesticks/ports/exchanges"
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
