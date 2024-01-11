package domain

import (
	candlesticks "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"
	"github.com/lerenn/cryptellation/svc/indicators/internal/app/ports/db"
)

type indicators struct {
	db           db.Port
	candlesticks candlesticks.Client
}

func New(db db.Port, candlesticks candlesticks.Client) indicators {
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
