package nats

import (
	"context"
	"net/http"

	"github.com/lerenn/cryptellation/pkg/version"
	asyncapi "github.com/lerenn/cryptellation/svc/backtests/api/asyncapi"
	"github.com/lerenn/cryptellation/svc/backtests/internal/app"
)

type subscriber struct {
	backtests  app.Backtests
	controller *asyncapi.AppController
}

func newSubscriber(controller *asyncapi.AppController, app app.Backtests) subscriber {
	return subscriber{
		backtests:  app,
		controller: controller,
	}
}

func (s subscriber) AccountsListOperationReceived(ctx context.Context, msg asyncapi.AccountsListRequestMessage) error {
	return s.controller.ReplyToAccountsListOperation(ctx, msg, func(replyMsg *asyncapi.AccountsListResponseMessage) {
		// Get accounts
		accounts, err := s.backtests.GetAccounts(context.Background(), uint(msg.Payload.Id))
		if err != nil {
			replyMsg.Payload.Error = &asyncapi.ErrorSchema{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
			return
		}

		// Set message with accounts
		replyMsg.Set(accounts)
	})
}

func (s subscriber) AdvanceOperationReceived(ctx context.Context, msg asyncapi.AdvanceRequestMessage) error {
	return s.controller.ReplyToAdvanceOperation(ctx, msg, func(replyMsg *asyncapi.AdvanceResponseMessage) {
		err := s.backtests.Advance(context.Background(), uint(msg.Payload.Id))
		if err != nil {
			replyMsg.Payload.Error = &asyncapi.ErrorSchema{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
	})
}

func (s subscriber) CreateOperationReceived(ctx context.Context, msg asyncapi.CreateRequestMessage) error {
	return s.controller.ReplyToCreateOperation(ctx, msg, func(replyMsg *asyncapi.CreateResponseMessage) {
		// Get model request from message payload
		req, err := msg.ToModel()
		if err != nil {
			replyMsg.Payload.Error = &asyncapi.ErrorSchema{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}
			return
		}

		// Create backtest
		id, err := s.backtests.Create(context.Background(), req)
		if err != nil {
			replyMsg.Payload.Error = &asyncapi.ErrorSchema{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
			return
		}

		// Set response ID
		replyMsg.Payload.Id = int64(id)
	})
}

func (s subscriber) OrdersCreateOperationReceived(ctx context.Context, msg asyncapi.OrdersCreateRequestMessage) error {
	return s.controller.ReplyToOrdersCreateOperation(ctx, msg, func(replyMsg *asyncapi.OrdersCreateResponseMessage) {
		// Set order model from API
		order, err := msg.ToModel()
		if err != nil {
			replyMsg.Payload.Error = &asyncapi.ErrorSchema{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}
			return
		}

		// Create the order
		err = s.backtests.CreateOrder(context.Background(), uint(msg.Payload.Id), order)
		if err != nil {
			replyMsg.Payload.Error = &asyncapi.ErrorSchema{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
			return
		}
	})
}

func (s subscriber) OrdersListOperationReceived(ctx context.Context, msg asyncapi.OrdersListRequestMessage) error {
	return s.controller.ReplyToOrdersListOperation(ctx, msg, func(replyMsg *asyncapi.OrdersListResponseMessage) {
		// Get list of orders
		list, err := s.backtests.GetOrders(context.Background(), uint(msg.Payload.Id))
		if err != nil {
			replyMsg.Payload.Error = &asyncapi.ErrorSchema{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
			return
		}

		// Set message
		replyMsg.Set(list)
	})
}

func (s subscriber) SubscribeOperationReceived(ctx context.Context, msg asyncapi.SubscribeRequestMessage) error {
	return s.controller.ReplyToSubscribeOperation(ctx, msg, func(replyMsg *asyncapi.SubscribeResponseMessage) {
		err := s.backtests.SubscribeToEvents(
			context.Background(),
			uint(msg.Payload.Id),
			string(msg.Payload.Exchange),
			string(msg.Payload.Pair),
		)
		if err != nil {
			replyMsg.Payload.Error = &asyncapi.ErrorSchema{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
			return
		}
	})
}

func (s subscriber) ServiceInfoOperationReceived(ctx context.Context, msg asyncapi.ServiceInfoRequestMessage) error {
	return s.controller.ReplyToServiceInfoOperation(ctx, msg, func(replyMsg *asyncapi.ServiceInfoResponseMessage) {
		replyMsg.Payload.ApiVersion = asyncapi.AsyncAPIVersion
		replyMsg.Payload.BinVersion = version.Version()
	})
}
