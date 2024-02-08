package nats

import (
	"context"
	"fmt"
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

func (s subscriber) RegisterToTicksRequest(ctx context.Context, msg asyncapi.RegisteringRequestMessage) {
	telemetry.L(ctx).Info(fmt.Sprintf("Received register request: %+v\n", msg))

	// Set response
	resp := asyncapi.NewRegisteringResponseMessage()
	resp.SetAsResponseFrom(&msg)
	defer func() { _ = s.controller.PublishRegisterToTicksResponse(ctx, resp) }()

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
}

func (s subscriber) UnregisterToTicksRequest(ctx context.Context, msg asyncapi.RegisteringRequestMessage) {
	telemetry.L(ctx).Info(fmt.Sprintf("Received unregister request: %+v\n", msg))

	// Set response
	resp := asyncapi.NewRegisteringResponseMessage()
	resp.SetAsResponseFrom(&msg)
	defer func() { _ = s.controller.PublishUnregisterToTicksResponse(ctx, resp) }()

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
