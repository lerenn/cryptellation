// Generate code for grpc
//go:generate go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
//go:generate go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
//go:generate protoc --proto_path=../../api --go_out=./proto --go_opt=paths=source_relative --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=./proto --go-grpc_opt=paths=source_relative ticks.proto

package client

import (
	"context"
	"fmt"
	"os"

	"github.com/digital-feather/cryptellation/services/ticks/clients/go/proto"
	"github.com/digital-feather/cryptellation/services/ticks/internal/application/ports/pubsub"
	"github.com/digital-feather/cryptellation/services/ticks/internal/infrastructure/pubsub/nats"
	"github.com/digital-feather/cryptellation/services/ticks/pkg/models/tick"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	grpcClient proto.TicksServiceClient
	psClient   pubsub.Adapter
}

func New() (client *Client, close func() error, err error) {
	grpcAddr := os.Getenv("CRYPTELLATION_TICKS_GRPC_URL")
	if grpcAddr == "" {
		return nil, func() error { return nil }, xerrors.New("no grpc url provided")
	}

	conn, err := grpc.Dial(grpcAddr, grpcDialOpts(grpcAddr)...)
	if err != nil {
		return nil, func() error { return nil }, fmt.Errorf("dialing ticks grpc server: %w", err)
	}

	natsClient, err := nats.New()
	if err != nil {
		return nil, conn.Close, fmt.Errorf("creating NATs Client: %w", err)
	}

	return &Client{
			grpcClient: proto.NewTicksServiceClient(conn),
			psClient:   natsClient,
		}, func() error {
			natsClient.Close()
			return conn.Close()
		}, nil
}

func (c *Client) Register(ctx context.Context, exchange, symbol string) error {
	_, err := c.grpcClient.Register(ctx, &proto.RegisterRequest{
		Exchange:   exchange,
		PairSymbol: symbol,
	})
	return err
}

func (c *Client) Unregister(ctx context.Context, exchange, symbol string) error {
	_, err := c.grpcClient.Register(ctx, &proto.RegisterRequest{
		Exchange:   exchange,
		PairSymbol: symbol,
	})
	return err
}

func (c *Client) Listen(symbol string) (<-chan tick.Tick, error) {
	return c.psClient.Subscribe(symbol)
}

func grpcDialOpts(grpcAddr string) []grpc.DialOption {
	return []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
}
