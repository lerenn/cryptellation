package nats

import (
	"context"

	"cryptellation/internal/adapters/telemetry"
	"cryptellation/pkg/version"

	"cryptellation/svc/ticks/api/asyncapi"
	"cryptellation/svc/ticks/internal/app"
)

type subscriber struct {
	ticks      app.Ticks
	controller *asyncapi.AppController
}

func newSubscriber(controller *asyncapi.AppController, app app.Ticks) subscriber {
	return subscriber{
		ticks:      app,
		controller: controller,
	}
}

func (sub subscriber) ListeningOperationReceived(ctx context.Context, msg asyncapi.ListeningNotificationMessage) error {
	sub.ticks.ListeningNotificationReceived(ctx, msg.ToModel())
	return nil
}

func (s subscriber) ServiceInfoOperationReceived(ctx context.Context, msg asyncapi.ServiceInfoRequestMessage) error {
	return s.controller.ReplyToServiceInfoOperation(ctx, msg, func(replyMsg *asyncapi.ServiceInfoResponseMessage) {
		replyMsg.Payload.ApiVersion = asyncapi.AsyncAPIVersion
		replyMsg.Payload.BinVersion = version.Version()
		telemetry.L(ctx).Debugf("Request for service info received, replying with %+v", replyMsg.Payload)
	})
}
