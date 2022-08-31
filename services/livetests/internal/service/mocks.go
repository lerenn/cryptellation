package service

import (
	"context"

	"github.com/digital-feather/cryptellation/services/ticks/pkg/models/tick"
)

type mockedTicksClient struct {
}

func (c *mockedTicksClient) Register(ctx context.Context, exchange, symbol string) error {
	return nil
}

func (c *mockedTicksClient) Unregister(ctx context.Context, exchange, symbol string) error {
	return nil
}

func (c *mockedTicksClient) Listen(symbol string) (<-chan tick.Tick, error) {
	return nil, nil
}
