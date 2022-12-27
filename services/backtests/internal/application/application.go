package application

import (
	"github.com/digital-feather/cryptellation/services/backtests/internal/application/operations/backtests"
	"github.com/digital-feather/cryptellation/services/backtests/internal/application/ports/db"
	"github.com/digital-feather/cryptellation/services/backtests/internal/application/ports/pubsub"
	candlesticks "github.com/digital-feather/cryptellation/services/candlesticks/pkg/client"
)

type Application struct {
	Backtests backtests.Operator
}

func New(cs candlesticks.Interfacer, db db.Adapter, ps pubsub.Adapter) (*Application, error) {
	return &Application{
		Backtests: backtests.New(db, ps, cs),
	}, nil
}
