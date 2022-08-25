package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	app "github.com/digital-feather/cryptellation/services/ticks/internal/application"
	"github.com/digital-feather/cryptellation/services/ticks/pkg/client/proto"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
)

type GrpcController struct {
	application app.Application
	server      *grpc.Server
}

func New(application app.Application) GrpcController {
	return GrpcController{application: application}
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
	proto.RegisterTicksServiceServer(grpcServer, g)

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

func (g GrpcController) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	count, err := g.application.Commands.RegisterSymbolListener.Handle(ctx, req.Exchange, req.PairSymbol)
	if err != nil {
		return nil, err
	}

	return &proto.RegisterResponse{
		RegisteredCount: count,
	}, nil
}

func (g GrpcController) Unregister(ctx context.Context, req *proto.UnregisterRequest) (*proto.UnregisterResponse, error) {
	count, err := g.application.Commands.UnregisterSymbolListener.Handle(ctx, req.Exchange, req.PairSymbol)
	if err != nil {
		return nil, err
	}

	return &proto.UnregisterResponse{
		RegisteredCount: count,
	}, nil
}
