package domain

import (
	candlesticks "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"
	"github.com/lerenn/cryptellation/svc/forwardtests/internal/app"
	"github.com/lerenn/cryptellation/svc/forwardtests/internal/app/ports/db"
)

// Test interface implementation
var _ app.Forwardtests = (*Forwardtests)(nil)

type Forwardtests struct {
	db           db.Port
	candlesticks candlesticks.Client
}

func New(db db.Port, csClient candlesticks.Client) *Forwardtests {
	if db == nil {
		panic("nil db")
	}

	if csClient == nil {
		panic("nil candlesticks client")
	}

	return &Forwardtests{
		db:           db,
		candlesticks: csClient,
	}
}
