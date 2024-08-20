package domain

import (
	candlesticks "github.com/lerenn/cryptellation/candlesticks/clients/go"

	"github.com/lerenn/cryptellation/forwardtests/internal/app"
	"github.com/lerenn/cryptellation/forwardtests/internal/app/ports/db"
)

// Test interface implementation
var _ app.ForwardTests = (*ForwardTests)(nil)

type ForwardTests struct {
	db           db.Port
	candlesticks candlesticks.Client
}

func New(db db.Port, csClient candlesticks.Client) *ForwardTests {
	if db == nil {
		panic("nil db")
	}

	if csClient == nil {
		panic("nil candlesticks client")
	}

	return &ForwardTests{
		db:           db,
		candlesticks: csClient,
	}
}
