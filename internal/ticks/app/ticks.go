package app

import (
	"github.com/digital-feather/cryptellation/internal/ticks/app/ports/db"
	"github.com/digital-feather/cryptellation/internal/ticks/app/ports/events"
	"github.com/digital-feather/cryptellation/internal/ticks/app/ports/exchanges"
)

type Ticks struct {
	db        db.Port
	events    events.Port
	exchanges exchanges.Port
}

func New(evts events.Port, db db.Port, exchanges exchanges.Port) *Ticks {
	if evts == nil {
		panic("nil events")
	}

	if db == nil {
		panic("nil vdb")
	}

	if exchanges == nil {
		panic("nil exchanges clients")
	}

	return &Ticks{
		events:    evts,
		exchanges: exchanges,
		db:        db,
	}
}
