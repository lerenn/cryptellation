package daemon

import (
	"github.com/lerenn/cryptellation/internal/controllers/nats"
	"github.com/lerenn/cryptellation/pkg/config"
)

type controllers struct {
	nats *nats.IndicatorsSubscriber
}

func newControllers(components components) (controllers, error) {
	nats, err := nats.NewIndicatorsSubscriber(config.LoadNATS(), components.indicators)
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
