package daemon

import (
	"cryptellation/internal/config"

	"cryptellation/svc/exchanges/internal/controllers/nats"
)

type controllers struct {
	nats *nats.Controller
}

func newControllers(components components) (controllers, error) {
	nats, err := nats.NewController(config.LoadNATS(), components.exchanges)
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
