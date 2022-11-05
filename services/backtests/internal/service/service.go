package service

import (
	"log"

	"github.com/digital-feather/cryptellation/services/backtests/internal/application"
	candlesticksGrpc "github.com/digital-feather/cryptellation/services/candlesticks/pkg/client"
)

func NewApplication() (*application.Application, func(), error) {
	csClient, closeCsClient, err := candlesticksGrpc.New()
	if err != nil {
		return nil, func() {}, err
	}

	app, err := application.New(csClient)
	fn := func() {
		if err := closeCsClient(); err != nil {
			log.Println("error when closing candlestick client:", err)
		}
	}

	return app, fn, err
}

func NewMockedApplication() (*application.Application, error) {
	return application.New(MockedCandlesticksClient{})
}
