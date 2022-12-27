package application

import (
	"github.com/digital-feather/cryptellation/services/ticks/internal/application/operators/ticks"
	"github.com/digital-feather/cryptellation/services/ticks/internal/application/ports/db"
	"github.com/digital-feather/cryptellation/services/ticks/internal/application/ports/exchanges"
	"github.com/digital-feather/cryptellation/services/ticks/internal/application/ports/pubsub"
)

type Application struct {
	Ticks ticks.Operator
}

func New(db db.Adapter, ps pubsub.Adapter, exchanges map[string]exchanges.Adapter) (*Application, error) {
	return &Application{
		Ticks: ticks.New(ps, db, exchanges),
	}, nil
}
