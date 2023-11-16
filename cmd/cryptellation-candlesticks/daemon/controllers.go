package daemon

import (
	"github.com/lerenn/cryptellation/internal/controllers/nats"
	"github.com/lerenn/cryptellation/pkg/config"
)

type controllers struct {
	nats *nats.CandlesticksController
}

func newControllers(components components) (controllers, error) {
	nats, err := nats.NewCandlesticksController(config.LoadNATS(), components.candlesticks)
	if err != nil {
		return controllers{}, err
	}

	return controllers{
		nats: nats,
	}, nil
}

func (c controllers) AsyncListen() error {
	return c.nats.Listen()
}

func (c controllers) Close() {
	c.nats.Close()
}
