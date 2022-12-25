package backtests

import (
	"github.com/digital-feather/cryptellation/services/backtests/internal/application/ports/db"
	"github.com/digital-feather/cryptellation/services/backtests/internal/application/ports/pubsub"
	candlesticks "github.com/digital-feather/cryptellation/services/candlesticks/pkg/client"
)

// Test interface implementation
var _ Operator = (*Backtests)(nil)

type Backtests struct {
	db       db.Adapter
	pubsub   pubsub.Adapter
	csClient candlesticks.Client
}

func New(db db.Adapter, ps pubsub.Adapter, csClient candlesticks.Client) *Backtests {
	if db == nil {
		panic("nil db")
	}

	if ps == nil {
		panic("nil pubsub")
	}

	if csClient == nil {
		panic("nil candlesticks client")
	}

	return &Backtests{
		db:       db,
		pubsub:   ps,
		csClient: csClient,
	}
}
