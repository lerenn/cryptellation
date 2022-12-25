package application

import (
	"github.com/digital-feather/cryptellation/services/ticks/internal/application/operators/ticks"
	"github.com/digital-feather/cryptellation/services/ticks/internal/application/ports/exchanges"
	"github.com/digital-feather/cryptellation/services/ticks/internal/application/ports/pubsub"
	"github.com/digital-feather/cryptellation/services/ticks/internal/application/ports/vdb"
)

type Application struct {
	Ticks ticks.Operator
}

func New(db vdb.Adapter, ps pubsub.Adapter, exchanges map[string]exchanges.Adapter) (*Application, error) {
	return &Application{
		Ticks: ticks.New(ps, db, exchanges),
	}, nil
}
