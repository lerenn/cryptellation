package daemon

import (
	"github.com/lerenn/cryptellation/internal/controllers/nats"
	"github.com/lerenn/cryptellation/pkg/config"
)

type controllers struct {
	nats *nats.BacktestsController
}

func newControllers(components components) (controllers, error) {
	nats, err := nats.NewBacktestsController(config.LoadNATS(), components.backtests)
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
