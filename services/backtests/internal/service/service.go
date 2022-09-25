package service

import (
	"log"

	"github.com/digital-feather/cryptellation/services/backtests/internal/adapters/pubsub/nats"
	"github.com/digital-feather/cryptellation/services/backtests/internal/adapters/vdb/redis"
	app "github.com/digital-feather/cryptellation/services/backtests/internal/application"
	cmdBacktest "github.com/digital-feather/cryptellation/services/backtests/internal/application/commands/backtest"
	queriesBacktest "github.com/digital-feather/cryptellation/services/backtests/internal/application/queries/backtest"
	candlesticksClient "github.com/digital-feather/cryptellation/services/candlesticks/pkg/client"
	candlesticksGrpc "github.com/digital-feather/cryptellation/services/candlesticks/pkg/client"
)

func NewApplication() (app.Application, func(), error) {
	csClient, closeCsClient, err := candlesticksGrpc.New()
	if err != nil {
		return app.Application{}, func() {}, err
	}

	app, closeApp, err := newApplication(csClient)
	return app, func() {
		closeApp()
		if err := closeCsClient(); err != nil {
			log.Println("error when closing candlestick client:", err)
		}
	}, err
}

func NewMockedApplication() (app.Application, func(), error) {
	return newApplication(MockedCandlesticksClient{})
}

func newApplication(client candlesticksClient.Client) (app.Application, func(), error) {
	repository, err := redis.New()
	if err != nil {
		return app.Application{}, func() {}, err
	}

	ps, err := nats.New()
	if err != nil {
		return app.Application{}, func() {}, err
	}

	return app.Application{
		Commands: app.Commands{
			Backtest: app.BacktestCommands{
				Advance:           cmdBacktest.NewAdvanceHandler(repository, ps, client),
				Create:            cmdBacktest.NewCreateHandler(repository),
				CreateOrder:       cmdBacktest.NewCreateOrderHandler(repository, client),
				SubscribeToEvents: cmdBacktest.NewSubscribeToEventsHandler(repository),
			},
		},
		Queries: app.Queries{
			Backtest: app.BacktestQueries{
				GetAccounts: queriesBacktest.NewGetAccounts(repository),
				GetOrders:   queriesBacktest.NewGetOrders(repository),
			},
		},
	}, func() {}, nil
}
