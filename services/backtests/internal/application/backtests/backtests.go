package backtests

import (
	"github.com/digital-feather/cryptellation/services/backtests/internal/adapters/db"
	"github.com/digital-feather/cryptellation/services/backtests/internal/adapters/pubsub"
	candlesticks "github.com/digital-feather/cryptellation/services/candlesticks/pkg/client"
)

// Test interface implementation
var _ Operator = (*Backtests)(nil)

type Backtests struct {
	repository db.Adapter
	pubsub     pubsub.Adapter
	csClient   candlesticks.Client
}

func New(repository db.Adapter, ps pubsub.Adapter, csClient candlesticks.Client) *Backtests {
	if repository == nil {
		panic("nil repository")
	}

	if ps == nil {
		panic("nil pubsub")
	}

	if csClient == nil {
		panic("nil candlesticks client")
	}

	return &Backtests{
		repository: repository,
		pubsub:     ps,
		csClient:   csClient,
	}
}
