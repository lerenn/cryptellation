package nats

import (
	"context"
	"net/http"

	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
	"github.com/lerenn/cryptellation/pkg/version"

	asyncapi "github.com/lerenn/cryptellation/candlesticks/api/asyncapi"
	"github.com/lerenn/cryptellation/candlesticks/internal/app"
)

type candlesticksSubscriber struct {
	candlesticks app.Candlesticks
	controller   *asyncapi.AppController
}

func newCandlesticksSubscriber(controller *asyncapi.AppController, app app.Candlesticks) candlesticksSubscriber {
	return candlesticksSubscriber{
		candlesticks: app,
		controller:   controller,
	}
}

func (s candlesticksSubscriber) ListOperationReceived(ctx context.Context, msg asyncapi.ListRequestMessage) error {
	return s.controller.ReplyToListOperation(ctx, msg, func(replyMsg *asyncapi.ListResponseMessage) {
		// Set list payload
		payload, err := msg.ToModel()
		if err != nil {
			replyMsg.Payload.Error = &asyncapi.ErrorSchema{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}
			return
		}

		// Request list
		list, err := s.candlesticks.GetCached(context.Background(), payload)
		if err != nil {
			replyMsg.Payload.Error = &asyncapi.ErrorSchema{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
			return
		}

		// Add list to response
		if err := replyMsg.Set(list); err != nil {
			replyMsg.Payload.Error = &asyncapi.ErrorSchema{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
			return
		}
	})
}

func (s candlesticksSubscriber) ServiceInfoOperationReceived(ctx context.Context, msg asyncapi.ServiceInfoRequestMessage) error {
	return s.controller.ReplyToServiceInfoOperation(ctx, msg, func(replyMsg *asyncapi.ServiceInfoResponseMessage) {
		replyMsg.Payload.ApiVersion = asyncapi.AsyncAPIVersion
		replyMsg.Payload.BinVersion = version.Version()
		telemetry.L(ctx).Debugf("Request for service info received, replying with %+v", replyMsg.Payload)
	})
}
