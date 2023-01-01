// Generate code for grpc
//go:generate go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
//go:generate go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
//go:generate protoc --proto_path=../../api --go_out=./proto --go_opt=paths=source_relative --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=./proto --go-grpc_opt=paths=source_relative exchanges.proto

package client

import (
	"context"
	"fmt"
	"os"

	"github.com/digital-feather/cryptellation/services/exchanges/clients/go/proto"
	"github.com/digital-feather/cryptellation/services/exchanges/pkg/models/exchange"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	grpcClient proto.ExchangesServiceClient
}

func New() (client *Client, close func() error, err error) {
	grpcAddr := os.Getenv("CRYPTELLATION_EXCHANGES_GRPC_URL")
	if grpcAddr == "" {
		return nil, func() error { return nil }, xerrors.New("no grpc url provided")
	}

	conn, err := grpc.Dial(grpcAddr, grpcDialOpts(grpcAddr)...)
	if err != nil {
		return nil, func() error { return nil }, fmt.Errorf("dialing exchanges grpc server: %w", err)
	}

	return &Client{
		grpcClient: proto.NewExchangesServiceClient(conn),
	}, conn.Close, nil
}

func (client *Client) ReadExchanges(ctx context.Context, names ...string) ([]exchange.Exchange, error) {
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
