package domain

import (
	"github.com/lerenn/cryptellation/svc/backtests/internal/app"
	"github.com/lerenn/cryptellation/svc/backtests/internal/app/ports/db"
	"github.com/lerenn/cryptellation/svc/backtests/internal/app/ports/events"

	candlesticks "github.com/lerenn/cryptellation/candlesticks/clients/go"
)

// Test interface implementation
var _ app.Backtests = (*Backtests)(nil)

type Backtests struct {
	db           db.Port
	events       events.Port
	candlesticks candlesticks.Client
}

func New(db db.Port, evts events.Port, csClient candlesticks.Client) *Backtests {
	if db == nil {
		panic("nil db")
	}

	if evts == nil {
		panic("nil events")
	}

	if csClient == nil {
		panic("nil candlesticks client")
	}

	return &Backtests{
		db:           db,
		events:       evts,
		candlesticks: csClient,
	}
}
