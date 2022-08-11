package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/account"
	app "github.com/digital-feather/cryptellation/services/livetests/internal/application"
	"github.com/digital-feather/cryptellation/services/livetests/internal/domain/livetest"
	"github.com/digital-feather/cryptellation/services/livetests/pkg/client/proto"
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
	proto.RegisterLivetestsServiceServer(grpcServer, g)

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

func (g GrpcController) CreateLivetest(ctx context.Context, req *proto.CreateLivetestRequest) (*proto.CreateLivetestResponse, error) {
	newPayload, err := fromCreateLivetestRequest(req)
	if err != nil {
		return nil, err
	}

	id, err := g.application.Commands.Livetest.Create.Handle(ctx, newPayload)
	if err != nil {
		return nil, err
	}

	return &proto.CreateLivetestResponse{
		Id: uint64(id),
	}, nil
}

func fromCreateLivetestRequest(req *proto.CreateLivetestRequest) (livetest.NewPayload, error) {
	acc := make(map[string]account.Account, len(req.Accounts))
	for exch, v := range req.Accounts {
		balances := make(map[string]float64, len(v.Assets))
		for asset, qty := range v.Assets {
			balances[asset] = float64(qty)
		}

		acc[exch] = account.Account{
			Balances: balances,
		}
	}

	return livetest.NewPayload{
		Accounts: acc,
	}, nil
}
