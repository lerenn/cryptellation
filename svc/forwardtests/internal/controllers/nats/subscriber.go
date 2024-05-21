package nats

import (
	"context"

	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
	"github.com/lerenn/cryptellation/pkg/version"
	asyncapi "github.com/lerenn/cryptellation/svc/forwardtests/api/asyncapi"
	"github.com/lerenn/cryptellation/svc/forwardtests/internal/app"
)

type subscriber struct {
	forwardtests app.Forwardtests
	controller   *asyncapi.AppController
}

func newSubscriber(controller *asyncapi.AppController, app app.Forwardtests) subscriber {
	return subscriber{
		forwardtests: app,
		controller:   controller,
	}
}

func (s subscriber) ServiceInfoOperationReceived(ctx context.Context, msg asyncapi.ServiceInfoRequestMessage) error {
	return s.controller.ReplyToServiceInfoOperation(ctx, msg, func(replyMsg *asyncapi.ServiceInfoResponseMessage) {
		replyMsg.Payload.ApiVersion = asyncapi.AsyncAPIVersion
		replyMsg.Payload.BinVersion = version.Version()
		telemetry.L(ctx).Debugf("Request for service info received, replying with %+v", replyMsg.Payload)
	})
}
