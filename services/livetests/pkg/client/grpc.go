package client

import (
	"context"
	"fmt"
	"os"

	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/account"
	"github.com/digital-feather/cryptellation/services/livetests/pkg/client/proto"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	grpcClient proto.LivetestsServiceClient
}

func New() (client *GrpcClient, close func() error, err error) {
	grpcAddr := os.Getenv("CRYPTELLATION_LIVETESTS_GRPC_URL")
	if grpcAddr == "" {
		return nil, func() error { return nil }, xerrors.New("no grpc url provided")
	}

	conn, err := grpc.Dial(grpcAddr, grpcDialOpts(grpcAddr)...)
	if err != nil {
		return nil, func() error { return nil }, fmt.Errorf("dialing livetests grpc server: %w", err)
	}

	return &GrpcClient{
		grpcClient: proto.NewLivetestsServiceClient(conn),
	}, conn.Close, nil
}

func (c *GrpcClient) CreateLivetest(ctx context.Context, accounts map[string]account.Account) (id uint, err error) {
	pbAccounts := make(map[string]*proto.Account)
	for n, a := range accounts {
		pbAccounts[n] = accountToProtoBuf(a)
	}

	resp, err := c.grpcClient.CreateLivetest(ctx, &proto.CreateLivetestRequest{
		Accounts: pbAccounts,
	})

	return uint(resp.Id), err
}

func grpcDialOpts(grpcAddr string) []grpc.DialOption {
	return []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
}
