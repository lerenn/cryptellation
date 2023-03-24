package backtests

import (
	"context"
	"net/http"

	"github.com/digital-feather/cryptellation/services/backtests"
)

type subscriber struct {
	backtests  backtests.Interface
	controller *AppController
}

func newSubscriber(controller *AppController, app backtests.Interface) subscriber {
	return subscriber{
		backtests:  app,
		controller: controller,
	}
}

func (s subscriber) BacktestsAccountsListRequest(msg BacktestsAccountsListRequestMessage, _ bool) {
	// Prepare response and set send at the end
	resp := NewBacktestsAccountsListResponseMessage()
	resp.SetAsResponseFrom(msg)
	defer func() { _ = s.controller.PublishBacktestsAccountsListResponse(resp) }()

	// Get accounts
	accounts, err := s.backtests.GetAccounts(context.Background(), uint(msg.Payload.ID))
	if err != nil {
		resp.Payload.Error = &ErrorSchema{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}

	// Set message with accounts
	resp.Set(accounts)
}

func (s subscriber) BacktestsAdvanceRequest(msg BacktestsAdvanceRequestMessage, _ bool) {
	// Prepare response and set send at the end
	resp := NewBacktestsAdvanceResponseMessage()
	resp.SetAsResponseFrom(msg)
	defer func() { _ = s.controller.PublishBacktestsAdvanceResponse(resp) }()

	// Advance application
	err := s.backtests.Advance(context.Background(), uint(msg.Payload.ID))
	if err != nil {
		resp.Payload.Error = &ErrorSchema{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}
}

func (s subscriber) BacktestsCreateRequest(msg BacktestsCreateRequestMessage, _ bool) {
	// Prepare response and set send at the end
	resp := NewBacktestsCreateResponseMessage()
	resp.SetAsResponseFrom(msg)
	defer func() { _ = s.controller.PublishBacktestsCreateResponse(resp) }()

	// Get model request from message payload
	req, err := msg.ToModel()
	if err != nil {
		resp.Payload.Error = &ErrorSchema{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		return
	}

	// Create backtest
	id, err := s.backtests.Create(context.Background(), req)
	if err != nil {
		resp.Payload.Error = &ErrorSchema{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}

	// Set response ID
	resp.Payload.ID = int64(id)
}

func (s subscriber) BacktestsOrdersCreateRequest(msg BacktestsOrdersCreateRequestMessage, _ bool) {
	// Prepare response and set send at the end
	resp := NewBacktestsOrdersCreateResponseMessage()
	resp.SetAsResponseFrom(msg)
	defer func() { _ = s.controller.PublishBacktestsOrdersCreateResponse(resp) }()

	// Set order model from API
	order, err := msg.ToModel()
	if err != nil {
		resp.Payload.Error = &ErrorSchema{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		return
	}

	// Create the order
	err = s.backtests.CreateOrder(context.Background(), uint(msg.Payload.ID), order)
	if err != nil {
		resp.Payload.Error = &ErrorSchema{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}
}

func (s subscriber) BacktestsOrdersListRequest(msg BacktestsOrdersListRequestMessage, _ bool) {
	// Prepare response and set send at the end
	resp := NewBacktestsOrdersListResponseMessage()
	resp.SetAsResponseFrom(msg)
	defer func() { _ = s.controller.PublishBacktestsOrdersListResponse(resp) }()

	// Get list of orders
	list, err := s.backtests.GetOrders(context.Background(), uint(msg.Payload.ID))
	if err != nil {
		resp.Payload.Error = &ErrorSchema{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}

	// Set message
	resp.Set(list)
}

func (s subscriber) BacktestsSubscribeRequest(msg BacktestsSubscribeRequestMessage, _ bool) {
	// Prepare response and set send at the end
	resp := NewBacktestsSubscribeResponseMessage()
	resp.SetAsResponseFrom(msg)
	defer func() { _ = s.controller.PublishBacktestsSubscribeResponse(resp) }()

	// Set subscription
	err := s.backtests.SubscribeToEvents(
		context.Background(),
		uint(msg.Payload.ID),
		string(msg.Payload.ExchangeName),
		string(msg.Payload.PairSymbol),
	)
	if err != nil {
		resp.Payload.Error = &ErrorSchema{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}
}
