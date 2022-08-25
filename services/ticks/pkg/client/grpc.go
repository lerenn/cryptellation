package client

import (
	"context"
	"fmt"
	"os"

	"github.com/digital-feather/cryptellation/services/ticks/internal/adapters/pubsub"
	"github.com/digital-feather/cryptellation/services/ticks/internal/adapters/pubsub/nats"
	"github.com/digital-feather/cryptellation/services/ticks/pkg/client/proto"
	"github.com/digital-feather/cryptellation/services/ticks/pkg/models/tick"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	grpcClient proto.TicksServiceClient
	psClient   pubsub.Port
}

func New() (client *GrpcClient, close func() error, err error) {
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

	return &GrpcClient{
			grpcClient: proto.NewTicksServiceClient(conn),
			psClient:   natsClient,
		}, func() error {
			natsClient.Close()
			return conn.Close()
		}, nil
}

func (c *GrpcClient) Register(ctx context.Context, exchange, symbol string) error {
	_, err := c.grpcClient.Register(ctx, &proto.RegisterRequest{
		Exchange:   exchange,
		PairSymbol: symbol,
	})
	return err
}

func (c *GrpcClient) Unregister(ctx context.Context, exchange, symbol string) error {
	_, err := c.grpcClient.Register(ctx, &proto.RegisterRequest{
		Exchange:   exchange,
		PairSymbol: symbol,
	})
	return err
}

func (c *GrpcClient) Listen(symbol string) (<-chan tick.Tick, error) {
	return c.psClient.Subscribe(symbol)
}

func grpcDialOpts(grpcAddr string) []grpc.DialOption {
	return []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
}
