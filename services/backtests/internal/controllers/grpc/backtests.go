package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/digital-feather/cryptellation/services/backtests/clients/go/proto"
	"github.com/digital-feather/cryptellation/services/backtests/internal/application/domains/backtest"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/account"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/order"
)

func (g GrpcController) CreateBacktest(ctx context.Context, req *proto.CreateBacktestRequest) (*proto.CreateBacktestResponse, error) {
	newPayload, err := fromCreateBacktestRequest(req)
	if err != nil {
		return nil, err
	}

	id, err := g.app.Backtests.Create(ctx, newPayload)
	if err != nil {
		return nil, err
	}

	return &proto.CreateBacktestResponse{
		Id: uint64(id),
	}, nil
}

func fromCreateBacktestRequest(req *proto.CreateBacktestRequest) (backtest.NewPayload, error) {
	st, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return backtest.NewPayload{}, fmt.Errorf("error when parsing start_time: %w", err)
	}

	var et *time.Time
	if req.EndTime != "" {
		t, err := time.Parse(time.RFC3339, req.EndTime)
		if err != nil {
			return backtest.NewPayload{}, fmt.Errorf("error when parsing start_time: %w", err)
		}
		et = &t
	}

	var tbe *time.Duration
	if req.SecondsBetweenPriceEvents > 0 {
		d := time.Duration(req.SecondsBetweenPriceEvents) * time.Second
		tbe = &d
	}

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

	return backtest.NewPayload{
		Accounts:              acc,
		StartTime:             st,
		EndTime:               et,
		DurationBetweenEvents: tbe,
	}, nil
}

func (g GrpcController) AdvanceBacktest(ctx context.Context, req *proto.AdvanceBacktestRequest) (*proto.AdvanceBacktestResponse, error) {
	if err := g.app.Backtests.Advance(ctx, uint(req.Id)); err != nil {
		return nil, err
	}

	return &proto.AdvanceBacktestResponse{}, nil
}

func (g GrpcController) SubscribeToBacktestEvents(ctx context.Context, req *proto.SubscribeToBacktestEventsRequest) (*proto.SubscribeToBacktestEventsResponse, error) {
	err := g.app.Backtests.SubscribeToEvents(ctx, uint(req.Id), req.ExchangeName, req.PairSymbol)
	return &proto.SubscribeToBacktestEventsResponse{}, err
}

func (g GrpcController) CreateBacktestOrder(ctx context.Context, req *proto.CreateBacktestOrderRequest) (*proto.CreateBacktestOrderResponse, error) {
	order, err := order.FromProtoBuf(req.Order)
	if err != nil {
		return nil, err
	}

	err = g.app.Backtests.CreateOrder(ctx, uint(req.BacktestId), order)
	return &proto.CreateBacktestOrderResponse{}, err
}

func (g GrpcController) BacktestAccounts(ctx context.Context, req *proto.BacktestAccountsRequest) (*proto.BacktestAccountsResponse, error) {
	accounts, err := g.app.Backtests.GetAccounts(ctx, uint(req.BacktestId))
	if err != nil {
		return nil, err
	}

	resp := proto.BacktestAccountsResponse{
		Accounts: make(map[string]*proto.Account, len(accounts)),
	}

	for exch, acc := range accounts {
		resp.Accounts[exch] = acc.ToProtoBuf()
	}

	return &resp, nil
}

func (g GrpcController) BacktestOrders(ctx context.Context, req *proto.BacktestOrdersRequest) (*proto.BacktestOrdersResponse, error) {
	orders, err := g.app.Backtests.GetOrders(ctx, uint(req.BacktestId))
	if err != nil {
		return nil, err
	}

	formattedOrders := make([]*proto.Order, len(orders))
	for i, o := range orders {
		formattedOrders[i] = o.ToProtoBuf()
	}

	return &proto.BacktestOrdersResponse{
		Orders: formattedOrders,
	}, nil
}
