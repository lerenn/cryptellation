package nats

import (
	"context"
	"net/http"
	"time"

	"github.com/digital-feather/cryptellation/internal/backtests/app"
	"github.com/digital-feather/cryptellation/internal/backtests/app/domain"
	"github.com/digital-feather/cryptellation/internal/backtests/infra/events/nats/generated"
	"github.com/digital-feather/cryptellation/pkg/types/account"
	"github.com/digital-feather/cryptellation/pkg/types/period"
	"github.com/digital-feather/cryptellation/pkg/utils"
)

type subscriber struct {
	backtests  app.Controller
	controller *generated.AppController
}

func newSubscriber(controller *generated.AppController, app app.Controller) subscriber {
	return subscriber{
		backtests:  app,
		controller: controller,
	}
}

func (s subscriber) BacktestsAccountsListRequest(msg generated.BacktestsAccountsListRequestMessage, _ bool) {
	// Prepare response and set send at the end
	resp := generated.NewBacktestsAccountsListResponseMessage()
	resp.SetAsResponseFrom(msg)
	defer func() { _ = s.controller.PublishBacktestsAccountsListResponse(resp) }()

	// Get accounts
	accounts, err := s.backtests.GetAccounts(context.Background(), uint(msg.Payload.ID))
	if err != nil {
		resp.Payload.Error = &generated.ErrorSchema{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}

	// Format accounts
	respAccounts := make([]generated.AccountSchema, 0, len(accounts))
	for name, acc := range accounts {
		respAccounts = append(respAccounts, accountModelToAPI(name, acc))
	}

	// Set response
	resp.Payload.Accounts = respAccounts
}

func (s subscriber) BacktestsAdvanceRequest(msg generated.BacktestsAdvanceRequestMessage, _ bool) {
	// Prepare response and set send at the end
	resp := generated.NewBacktestsAdvanceResponseMessage()
	resp.SetAsResponseFrom(msg)
	defer func() { _ = s.controller.PublishBacktestsAdvanceResponse(resp) }()

	// Advance application
	err := s.backtests.Advance(context.Background(), uint(msg.Payload.ID))
	if err != nil {
		resp.Payload.Error = &generated.ErrorSchema{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}
}

func (s subscriber) BacktestsCreateRequest(msg generated.BacktestsCreateRequestMessage, _ bool) {
	// Prepare response and set send at the end
	resp := generated.NewBacktestsCreateResponseMessage()
	resp.SetAsResponseFrom(msg)
	defer func() { _ = s.controller.PublishBacktestsCreateResponse(resp) }()

	// Format accounts
	accounts := make(map[string]account.Account)
	for _, acc := range msg.Payload.Accounts {
		name, a := accountModelFromAPI(acc)
		accounts[name] = a
	}

	// Get duration between events
	var duration *time.Duration
	if msg.Payload.Period != nil {
		symbol, err := period.FromString((string)(*msg.Payload.Period))
		if err != nil {
			resp.Payload.Error = &generated.ErrorSchema{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}
			return
		}

		duration = utils.ToReference(symbol.Duration())
	}

	// Process message
	req := domain.NewPayload{
		Accounts:              accounts,
		StartTime:             time.Time(msg.Payload.StartTime),
		EndTime:               (*time.Time)(msg.Payload.EndTime),
		DurationBetweenEvents: duration,
	}

	// Create backtest
	id, err := s.backtests.Create(context.Background(), req)
	if err != nil {
		resp.Payload.Error = &generated.ErrorSchema{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}

	// Set response ID
	resp.Payload.ID = int64(id)
}

func (s subscriber) BacktestsOrdersCreateRequest(msg generated.BacktestsOrdersCreateRequestMessage, _ bool) {
	// Prepare response and set send at the end
	resp := generated.NewBacktestsOrdersCreateResponseMessage()
	resp.SetAsResponseFrom(msg)
	defer func() { _ = s.controller.PublishBacktestsOrdersCreateResponse(resp) }()

	// Set order model from API
	order, err := orderModelFromAPI(msg.Payload.Order)
	if err != nil {
		resp.Payload.Error = &generated.ErrorSchema{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		return
	}

	// Create the order
	err = s.backtests.CreateOrder(context.Background(), uint(msg.Payload.ID), order)
	if err != nil {
		resp.Payload.Error = &generated.ErrorSchema{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}
}

func (s subscriber) BacktestsOrdersListRequest(msg generated.BacktestsOrdersListRequestMessage, _ bool) {
	// Prepare response and set send at the end
	resp := generated.NewBacktestsOrdersListResponseMessage()
	resp.SetAsResponseFrom(msg)
	defer func() { _ = s.controller.PublishBacktestsOrdersListResponse(resp) }()

	// Get list of orders
	list, err := s.backtests.GetOrders(context.Background(), uint(msg.Payload.ID))
	if err != nil {
		resp.Payload.Error = &generated.ErrorSchema{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}

	// Return list
	resp.Payload.Orders = make([]generated.OrderSchema, len(list))
	for i, o := range list {
		resp.Payload.Orders[i] = orderModelToAPI(o)
	}
}

func (s subscriber) BacktestsSubscribeRequest(msg generated.BacktestsSubscribeRequestMessage, _ bool) {
	// Prepare response and set send at the end
	resp := generated.NewBacktestsSubscribeResponseMessage()
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
		resp.Payload.Error = &generated.ErrorSchema{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}
}
