package nats

import (
	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
	"github.com/lerenn/cryptellation/pkg/config"
)

type Adapter struct {
	broker extensions.BrokerController
}

func New(c config.NATS) (*Adapter, error) {
	// Validate configuration
	if err := c.Validate(); err != nil {
		return nil, err
	}

	// Create a NATS Controller
	broker, err := nats.NewController(c.URL())
	if err != nil {
		return nil, err
	}

	return &Adapter{
		broker: broker,
	}, nil
}

func (a Adapter) Broker() extensions.BrokerController {
	return a.broker
}
