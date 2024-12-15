package client

import (
	"github.com/lerenn/cryptellation/v1/clients/go/raw"
	"github.com/lerenn/cryptellation/v1/pkg/config"
)

type client struct {
	raw.Raw
}

// New creates a new client to execute temporal workflows.
// temporalConfig is the optional configuration to use for the temporal client.
func New(temporalConfig ...config.Temporal) (Client, error) {
	c, err := raw.New(temporalConfig...)
	if err != nil {
		return client{}, err
	}

	return client{
		Raw: c,
	}, nil
}

func (c client) RawClient() raw.Raw {
	return c.Raw
}
