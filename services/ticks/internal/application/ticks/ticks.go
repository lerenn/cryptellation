package ticks

import (
	"github.com/digital-feather/cryptellation/services/ticks/internal/adapters/exchanges"
	"github.com/digital-feather/cryptellation/services/ticks/internal/adapters/pubsub"
	"github.com/digital-feather/cryptellation/services/ticks/internal/adapters/vdb"
)

type Ticks struct {
	vdb       vdb.Adapter
	pubsub    pubsub.Adapter
	exchanges map[string]exchanges.Adapter
}

func New(ps pubsub.Adapter, db vdb.Adapter, exchanges map[string]exchanges.Adapter) *Ticks {
	if ps == nil {
		panic("nil pubsub")
	}

	if db == nil {
		panic("nil vdb")
	}

	if exchanges == nil {
		panic("nil exchanges clients")
	}

	return &Ticks{
		pubsub:    ps,
		exchanges: exchanges,
		vdb:       db,
	}
}
