package domain

import (
	candlesticks "cryptellation/svc/candlesticks/clients/go"

	"cryptellation/svc/indicators/internal/app/ports/db"
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
