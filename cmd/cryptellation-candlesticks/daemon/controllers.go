package daemon

import (
	"github.com/lerenn/cryptellation/internal/ctrl/candlesticks/events"
	"github.com/lerenn/cryptellation/pkg/config"
)

type controllers struct {
	nats *events.NATS
}

func newControllers(components components) (controllers, error) {
	nats, err := events.NewNATS(config.LoadNATSConfigFromEnv(), components.candlesticks)
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
