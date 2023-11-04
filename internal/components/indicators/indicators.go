package indicators

import (
	client "github.com/lerenn/cryptellation/clients/go"
	"github.com/lerenn/cryptellation/internal/components/indicators/ports/db"
)

type indicators struct {
	db           db.Port
	candlesticks client.Candlesticks
}

func New(db db.Port, candlesticks client.Candlesticks) indicators {
	if db == nil {
		panic("nil db")
	}

	if candlesticks == nil {
		panic("nil candlesticks")
	}

	return indicators{
		db:           db,
		candlesticks: candlesticks,
	}
}
