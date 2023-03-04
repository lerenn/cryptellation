package app

import (
	"github.com/digital-feather/cryptellation/internal/ticks/app/ports/db"
	"github.com/digital-feather/cryptellation/internal/ticks/app/ports/exchanges"
	"github.com/digital-feather/cryptellation/internal/ticks/app/ports/pubsub"
)

type Ticks struct {
	db        db.Port
	pubsub    pubsub.Port
	exchanges exchanges.Port
}

func New(ps pubsub.Port, db db.Port, exchanges exchanges.Port) *Ticks {
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
		db:        db,
	}
}
