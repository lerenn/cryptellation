package nats

import (
	"context"
	"net/http"

	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
	"github.com/lerenn/cryptellation/pkg/version"
	"github.com/lerenn/cryptellation/svc/ticks/api/asyncapi"
	ticks "github.com/lerenn/cryptellation/svc/ticks/internal/app"
)

type subscriber struct {
	ticks      ticks.Ticks
	controller *asyncapi.AppController
}

func newSubscriber(controller *asyncapi.AppController, app ticks.Ticks) subscriber {
	return subscriber{
		ticks:      app,
		controller: controller,
	}
}

func (s subscriber) RegisterOperationReceived(ctx context.Context, msg asyncapi.RegisteringRequestMessage) error {
	telemetry.L(ctx).Infof("Received register request: %+v\n", msg)

	return s.controller.ReplyToRegisterOperation(ctx, msg, func(resp *asyncapi.RegisteringResponseMessage) {
		// Register as requested
		count, err := s.ticks.Register(
			context.Background(),
			string(msg.Payload.Exchange),
			string(msg.Payload.Pair),
		)

		// If there is an error, return it as BadRequest
		// TODO: get correct error
		if err != nil {
			resp.Payload.Error = &asyncapi.ErrorSchema{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}
			return
		}

		// Otherwise, return count
		resp.Payload.Count = &count
	})
}

func (s subscriber) UnregisterOperationReceived(ctx context.Context, msg asyncapi.RegisteringRequestMessage) error {
	telemetry.L(ctx).Infof("Received unregister request: %+v\n", msg)

	return s.controller.ReplyToUnregisterOperation(ctx, msg, func(resp *asyncapi.RegisteringResponseMessage) {
		// Register as requested
		count, err := s.ticks.Unregister(
			context.Background(),
			string(msg.Payload.Exchange),
			string(msg.Payload.Pair),
		)

		// If there is an error, return it as BadRequest
		// TODO: get correct error
		if err != nil {
			resp.Payload.Error = &asyncapi.ErrorSchema{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}
			return
		}

		// Otherwise, return count
		resp.Payload.Count = &count
	})
}

func (s subscriber) ServiceInfoOperationReceived(ctx context.Context, msg asyncapi.ServiceInfoRequestMessage) error {
	return s.controller.ReplyToServiceInfoOperation(ctx, msg, func(replyMsg *asyncapi.ServiceInfoResponseMessage) {
		replyMsg.Payload.ApiVersion = asyncapi.AsyncAPIVersion
		replyMsg.Payload.BinVersion = version.Version()
	})
}
