package nats

//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.5.0 -g client -p gen -i ../../../../../api/asyncapi.yaml -o ./gen/client.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.5.0 -g broker -p gen -i ../../../../../api/asyncapi.yaml -o ./gen/broker.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.5.0 -g types  -p gen -i ../../../../../api/asyncapi.yaml -o ./gen/types.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.5.0 -g nats   -p gen -i ../../../../../api/asyncapi.yaml -o ./gen/nats.gen.go

import (
	"context"

	"github.com/digital-feather/cryptellation/services/candlesticks/clients/go/payloads"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/candlestick"
)

type Client struct {
}

func New() *Client {
	return &Client{}
}

func (c *Client) ReadCandlesticks(ctx context.Context, payload payloads.ReadCandlesticksPayload) (*candlestick.List, error) {
	// TODO
	return nil, nil
}
