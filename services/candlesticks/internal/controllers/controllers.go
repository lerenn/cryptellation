package controllers

import (
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/components/candlesticks"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/controllers/nats"
)

type Controller interface {
	Listen() error
	Close()
}

func New(c Config, candlesticks candlesticks.Port) ([]Controller, error) {
	natsSrv, err := nats.NewServer(c.NATS, candlesticks)
	if err != nil {
		return nil, err
	}

	return []Controller{
		natsSrv,
	}, nil
}
