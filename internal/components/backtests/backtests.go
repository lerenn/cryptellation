package backtests

import (
	client "github.com/lerenn/cryptellation/clients/go"
	"github.com/lerenn/cryptellation/internal/components/backtests/ports/db"
	"github.com/lerenn/cryptellation/internal/components/backtests/ports/events"
)

// Test interface implementation
var _ Interface = (*Backtests)(nil)

type Backtests struct {
	db           db.Port
	events       events.Port
	candlesticks client.Candlesticks
}

func New(db db.Port, evts events.Port, csClient client.Candlesticks) *Backtests {
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
