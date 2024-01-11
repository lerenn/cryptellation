package domain

import (
	"github.com/lerenn/cryptellation/svc/ticks/internal/app/ports/db"
	"github.com/lerenn/cryptellation/svc/ticks/internal/app/ports/events"
	"github.com/lerenn/cryptellation/svc/ticks/internal/app/ports/exchanges"
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
