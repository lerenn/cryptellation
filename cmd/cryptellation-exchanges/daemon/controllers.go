package daemon

import (
	"github.com/lerenn/cryptellation/internal/controllers/nats"
	"github.com/lerenn/cryptellation/pkg/config"
)

type controllers struct {
	nats *nats.ExchangesController
}

func newControllers(components components) (controllers, error) {
	nats, err := nats.NewExchangesController(config.LoadNATS(), components.exchanges)
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
