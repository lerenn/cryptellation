package client

import (
	"context"
	"fmt"
	"os"

	"github.com/digital-feather/cryptellation/services/exchanges/pkg/client/proto"
	"github.com/digital-feather/cryptellation/services/exchanges/pkg/models/exchange"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	grpcClient proto.ExchangesServiceClient
}

func New() (client *GrpcClient, close func() error, err error) {
	grpcAddr := os.Getenv("CRYPTELLATION_EXCHANGES_GRPC_URL")
	if grpcAddr == "" {
		return nil, func() error { return nil }, xerrors.New("no grpc url provided")
	}

	conn, err := grpc.Dial(grpcAddr, grpcDialOpts(grpcAddr)...)
	if err != nil {
		return nil, func() error { return nil }, fmt.Errorf("dialing exchanges grpc server: %w", err)
	}

	return &GrpcClient{
		grpcClient: proto.NewExchangesServiceClient(conn),
	}, conn.Close, nil
}

func (client *GrpcClient) ReadExchanges(ctx context.Context, names ...string) ([]exchange.Exchange, error) {
	resp, err := client.grpcClient.ReadExchanges(ctx, &proto.ReadExchangesRequest{
		Names: names,
	})
	if err != nil {
		return nil, err
	}

	exchanges := make([]exchange.Exchange, len(resp.Exchanges))
	for i, pbExch := range resp.Exchanges {
		exchanges[i], err = exchange.FromProtoBuf(pbExch)
		if err != nil {
			return nil, err
		}
	}

	return exchanges, nil
}

func grpcDialOpts(grpcAddr string) []grpc.DialOption {
	return []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
}
