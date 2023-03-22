package app

import (
	client "github.com/digital-feather/cryptellation/clients/go"
	"github.com/digital-feather/cryptellation/internal/backtests/app/ports/db"
	"github.com/digital-feather/cryptellation/internal/backtests/app/ports/events"
)

// Test interface implementation
var _ Controller = (*Backtests)(nil)

type Backtests struct {
	db           db.Adapter
	events       events.Adapter
	candlesticks client.Candlesticks
}

func New(db db.Adapter, evts events.Adapter, csClient client.Candlesticks) *Backtests {
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
