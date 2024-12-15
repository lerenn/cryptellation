package raw

import (
	"context"

	"github.com/lerenn/cryptellation/v1/pkg/config"
	temporalclient "go.temporal.io/sdk/client"
)

type raw struct {
	temporal temporalclient.Client
}

// New creates a new client to execute temporal workflows.
// temporalConfig is the optional configuration to use for the temporal client.
func New(temporalConfig ...config.Temporal) (Raw, error) {
	var t config.Temporal

	if len(temporalConfig) > 0 {
		t = temporalConfig[0]
	}

	// Load temporal configuration
	t = config.LoadTemporal(&t)
	if err := t.Validate(); err != nil {
		return raw{}, err
	}

	// Create temporal client
	c, err := t.CreateTemporalClient()
	if err != nil {
		return raw{}, err
	}

	return &raw{temporal: c}, nil
}

func (c raw) Temporal() temporalclient.Client {
	return c.temporal
}

// Close closes the client.
func (c raw) Close(_ context.Context) {
	c.temporal.Close()
}
