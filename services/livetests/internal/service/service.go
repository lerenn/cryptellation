package service

import (
	vdbRedis "github.com/digital-feather/cryptellation/services/livetests/internal/adapters/vdb/redis"
	app "github.com/digital-feather/cryptellation/services/livetests/internal/application"
	cmdLivetest "github.com/digital-feather/cryptellation/services/livetests/internal/application/commands/livetest"
	ticksClient "github.com/digital-feather/cryptellation/services/ticks/pkg/client"
)

func NewApplication() (app.Application, func() error, error) {
	tClient, tClientClose, err := ticksClient.New()
	if err != nil {
		return app.Application{}, func() error { return nil }, err
	}

	app, closeApp, err := newApplication(tClient)
	return app, func() error {
		closeApp()
		return tClientClose()
	}, err
}

func NewMockedApplication() (app.Application, func(), error) {
	return newApplication(&mockedTicksClient{})
}

func newApplication(tClient ticksClient.Client) (app.Application, func(), error) {
	repository, err := vdbRedis.New()
	if err != nil {
		return app.Application{}, func() {}, err
	}

	return app.Application{
		Commands: app.Commands{
			Livetest: app.LivetestCommands{
				Create: cmdLivetest.NewCreateHandler(repository, tClient),
			},
		},
		Queries: app.Queries{
			Livetest: app.LivetestQueries{},
		},
	}, func() {}, nil
}
