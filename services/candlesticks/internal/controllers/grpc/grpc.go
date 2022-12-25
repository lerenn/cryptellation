package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/digital-feather/cryptellation/services/candlesticks/internal/application"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/application/operators/candlesticks"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/client/proto"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/candlestick"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/period"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcController struct {
	app    *application.Application
	server *grpc.Server
}

func New(application *application.Application) GrpcController {
	return GrpcController{app: application}
}

func (g *GrpcController) Run() error {
	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		return xerrors.New("no service port provided")
	}
	addr := fmt.Sprintf(":%s", port)
	return g.RunOnAddr(addr)
}

func (g *GrpcController) RunOnAddr(addr string) error {
	grpcServer := grpc.NewServer()
	proto.RegisterCandlesticksServiceServer(grpcServer, g)

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("grpc listening error: %w", err)
	}

	log.Println("Starting: gRPC Listener")
	go func() {
		if err := grpcServer.Serve(listen); err != nil {
			log.Println("error when serving grpc:", err)
		}
	}()

	return nil
}

func (g *GrpcController) GracefulStop() {
	if g.server == nil {
		log.Println("WARNING: attempted to gracefully stop a non running grpc server")
		return
	}

	g.server.GracefulStop()
	g.server = nil
}

func (g *GrpcController) Stop() {
	if g.server == nil {
		log.Println("WARNING: attempted to stop a non running grpc server")
		return
	}

	g.server.Stop()
	g.server = nil
}

func (g GrpcController) ReadCandlesticks(ctx context.Context, req *proto.ReadCandlesticksRequest) (*proto.ReadCandlesticksResponse, error) {
	payload, err := fromReadCandlesticksRequest(req)
	if err != nil {
		return nil, err
	}

	list, err := g.app.Candlesticks.GetCached(ctx, payload)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	gcandlesticks := make([]*proto.Candlestick, 0, list.Len())
	_ = list.Loop(func(t time.Time, cs candlestick.Candlestick) (bool, error) {
		gcandlesticks = append(gcandlesticks, cs.ToProfoBuff(t))

		return false, nil
	})

	return &proto.ReadCandlesticksResponse{
		Candlesticks: gcandlesticks,
	}, nil
}

func fromReadCandlesticksRequest(req *proto.ReadCandlesticksRequest) (candlesticks.GetCachedPayload, error) {
	per, err := period.FromString(req.PeriodSymbol)
	if err != nil {
		return candlesticks.GetCachedPayload{}, err
	}

	payload := candlesticks.GetCachedPayload{
		ExchangeName: req.ExchangeName,
		PairSymbol:   req.PairSymbol,
		Period:       per,
		Limit:        uint(req.Limit),
	}

	if req.Start != "" {
		start, err := time.Parse(time.RFC3339Nano, req.Start)
		if err != nil {
			return candlesticks.GetCachedPayload{}, err
		}
		payload.Start = &start
	}

	if req.End != "" {
		end, err := time.Parse(time.RFC3339Nano, req.End)
		if err != nil {
			return candlesticks.GetCachedPayload{}, err
		}
		payload.End = &end
	}

	return payload, nil
}
