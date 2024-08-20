package nats

import (
	"context"
	"net/http"

	"cryptellation/pkg/adapters/telemetry"
	"cryptellation/pkg/version"

	asyncapi "cryptellation/svc/indicators/api/asyncapi"
	"cryptellation/svc/indicators/internal/app"
)

type subscriber struct {
	indicators app.Indicators
	controller *asyncapi.AppController
}

func newSubscriber(controller *asyncapi.AppController, app app.Indicators) subscriber {
	return subscriber{
		indicators: app,
		controller: controller,
	}
}

func (s subscriber) SMAOperationReceived(ctx context.Context, msg asyncapi.SMARequestMessage) error {
	return s.controller.ReplyToSMAOperation(ctx, msg, func(replyMsg *asyncapi.SMAResponseMessage) {
		// Change from requests type to application types
		payload, err := msg.ToModel()
		if err != nil {
			replyMsg.Payload.Error = &asyncapi.ErrorSchema{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}
			return
		}

		// Request exchange(s) information
		indicators, err := s.indicators.GetCachedSMA(context.Background(), payload)
		if err != nil {
			replyMsg.Payload.Error = &asyncapi.ErrorSchema{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
			return
		}

		// Add indicators to response
		replyMsg.Set(indicators)
	})
}

func (s subscriber) ServiceInfoOperationReceived(ctx context.Context, msg asyncapi.ServiceInfoRequestMessage) error {
	return s.controller.ReplyToServiceInfoOperation(ctx, msg, func(replyMsg *asyncapi.ServiceInfoResponseMessage) {
		replyMsg.Payload.ApiVersion = asyncapi.AsyncAPIVersion
		replyMsg.Payload.BinVersion = version.Version()
		telemetry.L(ctx).Debugf("Request for service info received, replying with %+v", replyMsg.Payload)
	})
}
