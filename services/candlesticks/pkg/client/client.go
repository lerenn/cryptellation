package client

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/client/proto"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/candlestick"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/period"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	grpcClient proto.CandlesticksServiceClient
}

func New() (client *Client, close func() error, err error) {
	grpcAddr := os.Getenv("CRYPTELLATION_CANDLESTICKS_GRPC_URL")
	if grpcAddr == "" {
		return nil, func() error { return nil }, xerrors.New("no grpc url provided")
	}

	conn, err := grpc.Dial(grpcAddr, grpcDialOpts(grpcAddr)...)
	if err != nil {
		return nil, func() error { return nil }, fmt.Errorf("dialing candlesticks grpc server: %w", err)
	}

	return &Client{
		grpcClient: proto.NewCandlesticksServiceClient(conn),
	}, conn.Close, nil
}

type ReadCandlestickPayload struct {
	ExchangeName string
	PairSymbol   string
	Period       period.Symbol
	Start        time.Time
	End          time.Time
	Limit        uint
}

func (c *Client) ReadCandlesticks(ctx context.Context, payload ReadCandlestickPayload) (*candlestick.List, error) {
	resp, err := c.grpcClient.ReadCandlesticks(ctx, &proto.ReadCandlesticksRequest{
		ExchangeName: payload.ExchangeName,
		PairSymbol:   payload.PairSymbol,
		PeriodSymbol: payload.Period.String(),
		Start:        payload.Start.Format(time.RFC3339Nano),
		End:          payload.End.Format(time.RFC3339Nano),
		Limit:        int64(payload.Limit),
	})
	if err != nil {
		return nil, err
	}

	l := candlestick.NewEmptyList(candlestick.ListID{
		ExchangeName: payload.ExchangeName,
		PairSymbol:   payload.PairSymbol,
		Period:       payload.Period,
	})

	return l, l.LoadFromProtoBuf(resp.Candlesticks)
}

func grpcDialOpts(grpcAddr string) []grpc.DialOption {
	return []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
}
