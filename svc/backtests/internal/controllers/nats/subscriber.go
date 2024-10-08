package nats

import (
	"context"
	"net/http"

	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
	"github.com/lerenn/cryptellation/pkg/version"

	asyncapi "github.com/lerenn/cryptellation/svc/backtests/api/asyncapi"
	"github.com/lerenn/cryptellation/svc/backtests/internal/app"

	"github.com/google/uuid"
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
		// Parse backtest ID
		id, err := uuid.Parse(string(msg.Payload.Id))
		if err != nil {
			replyMsg.Payload.Error = &asyncapi.ErrorSchema{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}
			return
		}

		// Get accounts
		accounts, err := s.backtests.GetAccounts(context.Background(), id)
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
		// Parse backtest ID
		id, err := uuid.Parse(string(msg.Payload.Id))
		if err != nil {
			replyMsg.Payload.Error = &asyncapi.ErrorSchema{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}
			return
		}

		// Advance backtest
		if err := s.backtests.Advance(context.Background(), id); err != nil {
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
		replyMsg.Payload.Id = id.String()
	})
}

func (s subscriber) GetOperationReceived(ctx context.Context, msg asyncapi.GetRequestMessage) error {
	return s.controller.ReplyToGetOperation(ctx, msg, func(replyMsg *asyncapi.GetResponseMessage) {
		// Parse backtest ID
		id, err := uuid.Parse(string(msg.Payload.Id))
		if err != nil {
			replyMsg.Payload.Error = &asyncapi.ErrorSchema{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}
			return
		}

		// Get backtest
		bt, err := s.backtests.Get(context.Background(), id)
		if err != nil {
			replyMsg.Payload.Error = &asyncapi.ErrorSchema{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
			return
		}

		// Set message
		replyMsg.Set(bt)
	})
}

func (s subscriber) OrdersCreateOperationReceived(ctx context.Context, msg asyncapi.OrdersCreateRequestMessage) error {
	return s.controller.ReplyToOrdersCreateOperation(ctx, msg, func(replyMsg *asyncapi.OrdersCreateResponseMessage) {
		// Parse backtest ID
		id, err := uuid.Parse(string(msg.Payload.Id))
		if err != nil {
			replyMsg.Payload.Error = &asyncapi.ErrorSchema{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}
			return
		}

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
		if err := s.backtests.CreateOrder(context.Background(), id, order); err != nil {
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
		// Parse backtest ID
		id, err := uuid.Parse(string(msg.Payload.Id))
		if err != nil {
			replyMsg.Payload.Error = &asyncapi.ErrorSchema{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}
			return
		}

		// Get list of orders
		list, err := s.backtests.GetOrders(context.Background(), id)
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

func (s subscriber) ListOperationReceived(ctx context.Context, msg asyncapi.ListRequestMessage) error {
	return s.controller.ReplyToListOperation(ctx, msg, func(replyMsg *asyncapi.ListResponseMessage) {
		// Get list of backtests
		list, err := s.backtests.List(context.Background())
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
		// Parse backtest ID
		id, err := uuid.Parse(string(msg.Payload.Id))
		if err != nil {
			replyMsg.Payload.Error = &asyncapi.ErrorSchema{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}
			return
		}

		// Subscribe to events
		if err := s.backtests.SubscribeToEvents(
			context.Background(),
			id,
			string(msg.Payload.Exchange),
			string(msg.Payload.Pair),
		); err != nil {
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
		telemetry.L(ctx).Debugf("Request for service info received, replying with %+v", replyMsg.Payload)
	})
}
