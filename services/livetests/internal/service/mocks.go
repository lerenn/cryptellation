package service

import (
	"context"
	"time"

	"google.golang.org/grpc"

	ticksProto "github.com/digital-feather/cryptellation/services/ticks/pkg/client/proto"
)

type MockedTicksClient struct {
}

func (m MockedTicksClient) ListenSymbol(
	ctx context.Context,
	in *ticksProto.ListenSymbolRequest,
	opts ...grpc.CallOption,
) (ticksProto.TicksService_ListenSymbolClient, error) {
	var count int64
	return MockListenSymbolClient{
		exchange:   in.Exchange,
		pairSymbol: in.PairSymbol,
		count:      &count,
	}, nil
}

type MockListenSymbolClient struct {
	exchange   string
	pairSymbol string
	count      *int64
	grpc.ClientStream
}

func (m MockListenSymbolClient) Recv() (*ticksProto.Tick, error) {
	t := &ticksProto.Tick{
		Time:       time.Unix(*m.count, 0).Format(time.RFC3339),
		Exchange:   m.exchange,
		PairSymbol: m.pairSymbol,
		Price:      float32(*m.count),
	}

	*m.count += 1
	return t, nil
}
