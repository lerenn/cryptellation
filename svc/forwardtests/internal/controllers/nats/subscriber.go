package nats

import (
	"context"
	"net/http"

	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
	"github.com/lerenn/cryptellation/pkg/version"

	asyncapi "github.com/lerenn/cryptellation/svc/forwardtests/api/asyncapi"
	"github.com/lerenn/cryptellation/svc/forwardtests/internal/app"

	"github.com/google/uuid"
)

type subscriber struct {
	forwardtests app.ForwardTests
	controller   *asyncapi.AppController
}

func newSubscriber(controller *asyncapi.AppController, app app.ForwardTests) subscriber {
	return subscriber{
		forwardtests: app,
		controller:   controller,
	}
}

func (s subscriber) AccountsListOperationReceived(ctx context.Context, msg asyncapi.AccountsListRequestMessage) error {
	return s.controller.ReplyToAccountsListOperation(ctx, msg, func(replyMsg *asyncapi.AccountsListResponseMessage) {
		// Parse forward test ID
		id, err := uuid.Parse(string(msg.Payload.Id))
		if err != nil {
			replyMsg.Payload.Error = &asyncapi.ErrorSchema{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}
			return
		}

		accounts, err := s.forwardtests.GetAccounts(ctx, id)
		if err != nil {
			replyMsg.Payload.Error = &asyncapi.ErrorSchema{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
			return
		}

		replyMsg.Set(accounts)
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

		id, err := s.forwardtests.Create(ctx, req)
		if err != nil {
			replyMsg.Payload.Error = &asyncapi.ErrorSchema{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
			return
		}

		replyMsg.Payload.Id = id.String()
	})
}

func (s subscriber) ListOperationReceived(ctx context.Context, msg asyncapi.ListRequestMessage) error {
	return s.controller.ReplyToListOperation(ctx, msg, func(replyMsg *asyncapi.ListResponseMessage) {
		list, err := s.forwardtests.List(ctx, app.ListFilters{})
		if err != nil {
			replyMsg.Payload.Error = &asyncapi.ErrorSchema{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
			return
		}

		replyMsg.Set(list)
	})
}

func (s subscriber) OrdersCreateOperationReceived(ctx context.Context, msg asyncapi.OrdersCreateRequestMessage) error {
	return s.controller.ReplyToOrdersCreateOperation(ctx, msg, func(replyMsg *asyncapi.OrdersCreateResponseMessage) {
		// Get model request from message payload
		req, err := msg.ToModel()
		if err != nil {
			replyMsg.Payload.Error = &asyncapi.ErrorSchema{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}
			return
		}

		id, err := uuid.Parse(string(msg.Payload.Id))
		if err != nil {
			replyMsg.Payload.Error = &asyncapi.ErrorSchema{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}
			return
		}

		if err := s.forwardtests.CreateOrder(ctx, id, req); err != nil {
			replyMsg.Payload.Error = &asyncapi.ErrorSchema{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
			return
		}
	})
}

func (s subscriber) GetStatusOperationReceived(ctx context.Context, msg asyncapi.StatusRequestMessage) error {
	return s.controller.ReplyToGetStatusOperation(ctx, msg, func(replyMsg *asyncapi.StatusResponseMessage) {
		id, err := uuid.Parse(string(msg.Payload.Id))
		if err != nil {
			replyMsg.Payload.Error = &asyncapi.ErrorSchema{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}
			return
		}

		status, err := s.forwardtests.GetStatus(ctx, id)
		if err != nil {
			replyMsg.Payload.Error = &asyncapi.ErrorSchema{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
			return
		}

		replyMsg.Set(status)
	})

}

func (s subscriber) ServiceInfoOperationReceived(ctx context.Context, msg asyncapi.ServiceInfoRequestMessage) error {
	return s.controller.ReplyToServiceInfoOperation(ctx, msg, func(replyMsg *asyncapi.ServiceInfoResponseMessage) {
		replyMsg.Payload.ApiVersion = asyncapi.AsyncAPIVersion
		replyMsg.Payload.BinVersion = version.Version()
		telemetry.L(ctx).Debugf("Request for service info received, replying with %+v", replyMsg.Payload)
	})
}
