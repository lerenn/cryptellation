// Generate code for grpc
//go:generate go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
//go:generate go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
//go:generate protoc --proto_path=../../api --go_out=./proto --go_opt=paths=source_relative --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=./proto --go-grpc_opt=paths=source_relative backtests.proto

package client

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/digital-feather/cryptellation/services/backtests/internal/application/ports/pubsub"
	"github.com/digital-feather/cryptellation/services/backtests/internal/infrastructure/pubsub/nats"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/client/proto"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/account"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/event"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/order"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	grpcClient proto.BacktestsServiceClient
	natsClient pubsub.Adapter
}

func New() (client *Client, close func() error, err error) {
	grpcAddr := os.Getenv("CRYPTELLATION_BACKTESTS_GRPC_URL")
	if grpcAddr == "" {
		return nil, func() error { return nil }, xerrors.New("no grpc url provided")
	}

	conn, err := grpc.Dial(grpcAddr, grpcDialOpts(grpcAddr)...)
	if err != nil {
		return nil, func() error { return nil }, fmt.Errorf("dialing backtests grpc server: %w", err)
	}

	natsClient, err := nats.New()
	if err != nil {
		return nil, conn.Close, fmt.Errorf("creating NATs Client: %w", err)
	}

	return &Client{
			grpcClient: proto.NewBacktestsServiceClient(conn),
			natsClient: natsClient,
		}, func() error {
			natsClient.Close()
			return conn.Close()
		}, nil
}

func (c *Client) CreateBacktest(ctx context.Context, start, end time.Time, accounts map[string]account.Account) (id uint64, err error) {
	pbAccounts := make(map[string]*proto.Account)
	for n, a := range accounts {
		pbAccounts[n] = a.ToProtoBuf()
	}

	resp, err := c.grpcClient.CreateBacktest(ctx, &proto.CreateBacktestRequest{
		StartTime: start.Format(time.RFC3339Nano),
		EndTime:   end.Format(time.RFC3339Nano),
		Accounts:  pbAccounts,
	})
	if err != nil {
		return 0, err
	}

	return resp.Id, nil
}

func (c *Client) AdvanceBacktest(ctx context.Context, backtestID uint64) error {
	_, err := c.grpcClient.AdvanceBacktest(ctx, &proto.AdvanceBacktestRequest{
		Id: backtestID,
	})

	return err
}

func (c *Client) BacktestAccounts(ctx context.Context, backtestID uint64) (map[string]account.Account, error) {
	resp, err := c.grpcClient.BacktestAccounts(ctx, &proto.BacktestAccountsRequest{
		BacktestId: backtestID,
	})

	if err != nil {
		return nil, err
	}

	accounts := make(map[string]account.Account)
	for n, a := range resp.Accounts {
		accounts[n] = account.FromProtoBuf(a)
	}

	return accounts, nil
}

func (c *Client) CreateBacktestOrder(ctx context.Context, backtestID uint64, o order.Order) error {
	_, err := c.grpcClient.CreateBacktestOrder(ctx, &proto.CreateBacktestOrderRequest{
		BacktestId: backtestID,
		Order:      o.ToProtoBuf(),
	})
	return err
}

func (c *Client) BacktestOrders(ctx context.Context, backtestID uint64) ([]order.Order, error) {
	resp, err := c.grpcClient.BacktestOrders(ctx, &proto.BacktestOrdersRequest{
		BacktestId: backtestID,
	})
	if err != nil {
		return nil, err
	}

	orders := make([]order.Order, len(resp.Orders))
	for i, pb := range resp.Orders {
		o, err := order.FromProtoBuf(pb)
		if err != nil {
			return nil, err
		}
		orders[i] = o
	}

	return orders, nil
}

func (c *Client) SubscribeToBacktestEvents(ctx context.Context, backtestID uint64, exchangeName, pairSymbol string) error {
	_, err := c.grpcClient.SubscribeToBacktestEvents(ctx, &proto.SubscribeToBacktestEventsRequest{
		Id:           backtestID,
		ExchangeName: exchangeName,
		PairSymbol:   pairSymbol,
	})
	return err
}

func (c *Client) ListenBacktest(backtestID uint) (<-chan event.Event, error) {
	return c.natsClient.Subscribe(backtestID)
}

func grpcDialOpts(grpcAddr string) []grpc.DialOption {
	return []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
}
