package service

import (
	"log"

	vdbRedis "github.com/digital-feather/cryptellation/services/livetests/internal/adapters/vdb/redis"
	app "github.com/digital-feather/cryptellation/services/livetests/internal/application"
	cmdLivetest "github.com/digital-feather/cryptellation/services/livetests/internal/application/commands/livetest"
	ticksGrpc "github.com/digital-feather/cryptellation/services/ticks/pkg/client"
	ticksProto "github.com/digital-feather/cryptellation/services/ticks/pkg/client/proto"
)

func NewApplication() (app.Application, func(), error) {
	ticksClient, closeTicksClient, err := ticksGrpc.New()
	if err != nil {
		return app.Application{}, func() {}, err
	}

	app, closeApp, err := newApplication(ticksClient)
	return app, func() {
		closeApp()
		if err := closeTicksClient(); err != nil {
			log.Println("error when closing ticks client:", err)
		}
	}, err
}

func NewMockedApplication() (app.Application, func(), error) {
	return newApplication(MockedTicksClient{})
}

func newApplication(client ticksProto.TicksServiceClient) (app.Application, func(), error) {
	repository, err := vdbRedis.New()
	if err != nil {
		return app.Application{}, func() {}, err
	}

	return app.Application{
		Commands: app.Commands{
			Livetest: app.LivetestCommands{
				Create: cmdLivetest.NewCreateHandler(repository),
			},
		},
		Queries: app.Queries{
			Livetest: app.LivetestQueries{},
		},
	}, func() {}, nil
}
