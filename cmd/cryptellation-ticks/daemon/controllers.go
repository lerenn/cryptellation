package daemon

import (
	"github.com/lerenn/cryptellation/internal/controllers/nats"
	"github.com/lerenn/cryptellation/pkg/config"
)

type controllers struct {
	nats *nats.TicksController
}

func newControllers(components components) (controllers, error) {
	nats, err := nats.NewTicksController(config.LoadNATS(), components.ticks)
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
