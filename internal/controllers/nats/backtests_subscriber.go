package nats

import (
	"context"
	"net/http"

	"github.com/lerenn/cryptellation/internal/components/backtests"
	asyncapi "github.com/lerenn/cryptellation/pkg/asyncapi/backtests"
)

type backtestsSubscriber struct {
	backtests  backtests.Interface
	controller *asyncapi.AppController
}

func newBacktestsSubscriber(controller *asyncapi.AppController, app backtests.Interface) backtestsSubscriber {
	return backtestsSubscriber{
		backtests:  app,
		controller: controller,
	}
}

func (s backtestsSubscriber) CryptellationBacktestsAccountsListRequest(ctx context.Context, msg asyncapi.BacktestsAccountsListRequestMessage) {
	// Prepare response and set send at the end
	resp := asyncapi.NewBacktestsAccountsListResponseMessage()
	resp.SetAsResponseFrom(&msg)
	defer func() { _ = s.controller.PublishCryptellationBacktestsAccountsListResponse(ctx, resp) }()

	// Get accounts
	accounts, err := s.backtests.GetAccounts(context.Background(), uint(msg.Payload.ID))
	if err != nil {
		resp.Payload.Error = &asyncapi.ErrorSchema{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}

	// Set message with accounts
	resp.Set(accounts)
}

func (s backtestsSubscriber) CryptellationBacktestsAdvanceRequest(ctx context.Context, msg asyncapi.BacktestsAdvanceRequestMessage) {
	// Prepare response and set send at the end
	resp := asyncapi.NewBacktestsAdvanceResponseMessage()
	resp.SetAsResponseFrom(&msg)
	defer func() { _ = s.controller.PublishCryptellationBacktestsAdvanceResponse(ctx, resp) }()

	// Advance application
	err := s.backtests.Advance(context.Background(), uint(msg.Payload.ID))
	if err != nil {
		resp.Payload.Error = &asyncapi.ErrorSchema{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}
}

func (s backtestsSubscriber) CryptellationBacktestsCreateRequest(ctx context.Context, msg asyncapi.BacktestsCreateRequestMessage) {
	// Prepare response and set send at the end
	resp := asyncapi.NewBacktestsCreateResponseMessage()
	resp.SetAsResponseFrom(&msg)
	defer func() { _ = s.controller.PublishCryptellationBacktestsCreateResponse(ctx, resp) }()

	// Get model request from message payload
	req, err := msg.ToModel()
	if err != nil {
		resp.Payload.Error = &asyncapi.ErrorSchema{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		return
	}

	// Create backtest
	id, err := s.backtests.Create(context.Background(), req)
	if err != nil {
		resp.Payload.Error = &asyncapi.ErrorSchema{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}

	// Set response ID
	resp.Payload.ID = int64(id)
}

func (s backtestsSubscriber) CryptellationBacktestsOrdersCreateRequest(ctx context.Context, msg asyncapi.BacktestsOrdersCreateRequestMessage) {
	// Prepare response and set send at the end
	resp := asyncapi.NewBacktestsOrdersCreateResponseMessage()
	resp.SetAsResponseFrom(&msg)
	defer func() { _ = s.controller.PublishCryptellationBacktestsOrdersCreateResponse(ctx, resp) }()

	// Set order model from API
	order, err := msg.ToModel()
	if err != nil {
		resp.Payload.Error = &asyncapi.ErrorSchema{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		return
	}

	// Create the order
	err = s.backtests.CreateOrder(context.Background(), uint(msg.Payload.ID), order)
	if err != nil {
		resp.Payload.Error = &asyncapi.ErrorSchema{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}
}

func (s backtestsSubscriber) CryptellationBacktestsOrdersListRequest(ctx context.Context, msg asyncapi.BacktestsOrdersListRequestMessage) {
	// Prepare response and set send at the end
	resp := asyncapi.NewBacktestsOrdersListResponseMessage()
	resp.SetAsResponseFrom(&msg)
	defer func() { _ = s.controller.PublishCryptellationBacktestsOrdersListResponse(ctx, resp) }()

	// Get list of orders
	list, err := s.backtests.GetOrders(context.Background(), uint(msg.Payload.ID))
	if err != nil {
		resp.Payload.Error = &asyncapi.ErrorSchema{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}

	// Set message
	resp.Set(list)
}

func (s backtestsSubscriber) CryptellationBacktestsSubscribeRequest(ctx context.Context, msg asyncapi.BacktestsSubscribeRequestMessage) {
	// Prepare response and set send at the end
	resp := asyncapi.NewBacktestsSubscribeResponseMessage()
	resp.SetAsResponseFrom(&msg)
	defer func() { _ = s.controller.PublishCryptellationBacktestsSubscribeResponse(ctx, resp) }()

	// Set subscription
	err := s.backtests.SubscribeToEvents(
		context.Background(),
		uint(msg.Payload.ID),
		string(msg.Payload.ExchangeName),
		string(msg.Payload.PairSymbol),
	)
	if err != nil {
		resp.Payload.Error = &asyncapi.ErrorSchema{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}
}
