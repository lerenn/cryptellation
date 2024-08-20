package nats

import (
	"context"
	"net/http"

	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
	"github.com/lerenn/cryptellation/pkg/version"

	asyncapi "github.com/lerenn/cryptellation/svc/exchanges/api/asyncapi"
	exchanges "github.com/lerenn/cryptellation/svc/exchanges/internal/app"
)

type subscriber struct {
	exchanges  exchanges.Exchanges
	controller *asyncapi.AppController
}

func newSubscriber(controller *asyncapi.AppController, app exchanges.Exchanges) subscriber {
	return subscriber{
		exchanges:  app,
		controller: controller,
	}
}

func (s subscriber) ListOperationReceived(ctx context.Context, msg asyncapi.ListRequestMessage) error {
	return s.controller.ReplyToListOperation(ctx, msg, func(replyMsg *asyncapi.ListResponseMessage) {
		// Change from requests type to application types
		exchangesNames := msg.ToModel()

		// Request exchange(s) information
		exchanges, err := s.exchanges.GetCached(context.Background(), exchangesNames...)
		if err != nil {
			replyMsg.Payload.Error = &asyncapi.ErrorSchema{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
			return
		}

		// Add exchanges to response
		replyMsg.Set(exchanges)
	})
}

func (s subscriber) ServiceInfoOperationReceived(ctx context.Context, msg asyncapi.ServiceInfoRequestMessage) error {
	return s.controller.ReplyToServiceInfoOperation(ctx, msg, func(replyMsg *asyncapi.ServiceInfoResponseMessage) {
		replyMsg.Payload.ApiVersion = asyncapi.AsyncAPIVersion
		replyMsg.Payload.BinVersion = version.Version()
		telemetry.L(ctx).Debugf("Request for service info received, replying with %+v", replyMsg.Payload)
	})
}
