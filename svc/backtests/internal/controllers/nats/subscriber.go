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

func (s subscriber) ListBacktestAccountsRequest(ctx context.Context, msg asyncapi.ListBacktestAccountsRequestMessage) {
	// Prepare response and set send at the end
	resp := asyncapi.NewListBacktestAccountsResponseMessage()
	resp.SetAsResponseFrom(&msg)
	defer func() { _ = s.controller.PublishListBacktestAccountsResponse(ctx, resp) }()

	// Get accounts
	accounts, err := s.backtests.GetAccounts(context.Background(), uint(msg.Payload.Id))
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

func (s subscriber) AdvanceBacktestRequest(ctx context.Context, msg asyncapi.AdvanceBacktestRequestMessage) {
	// Prepare response and set send at the end
	resp := asyncapi.NewAdvanceBacktestResponseMessage()
	resp.SetAsResponseFrom(&msg)
	defer func() { _ = s.controller.PublishAdvanceBacktestResponse(ctx, resp) }()

	// Advance application
	err := s.backtests.Advance(context.Background(), uint(msg.Payload.Id))
	if err != nil {
		resp.Payload.Error = &asyncapi.ErrorSchema{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}
}

func (s subscriber) CreateBacktestRequest(ctx context.Context, msg asyncapi.CreateBacktestRequestMessage) {
	// Prepare response and set send at the end
	resp := asyncapi.NewCreateBacktestResponseMessage()
	resp.SetAsResponseFrom(&msg)
	defer func() { _ = s.controller.PublishCreateBacktestResponse(ctx, resp) }()

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
	resp.Payload.Id = int64(id)
}

func (s subscriber) CreateBacktestOrderRequest(ctx context.Context, msg asyncapi.CreateBacktestOrderRequestMessage) {
	// Prepare response and set send at the end
	resp := asyncapi.NewCreateBacktestOrderResponseMessage()
	resp.SetAsResponseFrom(&msg)
	defer func() { _ = s.controller.PublishCreateBacktestOrderResponse(ctx, resp) }()

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
	err = s.backtests.CreateOrder(context.Background(), uint(msg.Payload.Id), order)
	if err != nil {
		resp.Payload.Error = &asyncapi.ErrorSchema{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}
}

func (s subscriber) ListBacktestOrdersRequest(ctx context.Context, msg asyncapi.ListBacktestOrdersRequestMessage) {
	// Prepare response and set send at the end
	resp := asyncapi.NewListBacktestOrdersResponseMessage()
	resp.SetAsResponseFrom(&msg)
	defer func() { _ = s.controller.PublishListBacktestOrdersResponse(ctx, resp) }()

	// Get list of orders
	list, err := s.backtests.GetOrders(context.Background(), uint(msg.Payload.Id))
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

func (s subscriber) SubscribeBacktestRequest(ctx context.Context, msg asyncapi.SubscribeBacktestRequestMessage) {
	// Prepare response and set send at the end
	resp := asyncapi.NewSubscribeBacktestResponseMessage()
	resp.SetAsResponseFrom(&msg)
	defer func() { _ = s.controller.PublishSubscribeBacktestResponse(ctx, resp) }()

	// Set subscription
	err := s.backtests.SubscribeToEvents(
		context.Background(),
		uint(msg.Payload.Id),
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

func (s subscriber) ServiceInfoRequest(ctx context.Context, msg asyncapi.ServiceInfoRequestMessage) {
	// Prepare response and set send at the end
	resp := asyncapi.NewServiceInfoResponseMessage()
	resp.SetAsResponseFrom(&msg)
	defer func() { _ = s.controller.PublishServiceInfoResponse(ctx, resp) }()

	// Set info
	resp.Payload.ApiVersion = asyncapi.AsyncAPIVersion
	resp.Payload.BinVersion = version.Version()
}
