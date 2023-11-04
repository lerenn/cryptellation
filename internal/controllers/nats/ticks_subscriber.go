package nats

import (
	"context"
	"log"
	"net/http"

	"github.com/lerenn/cryptellation/internal/components/ticks"
	asyncapi "github.com/lerenn/cryptellation/pkg/asyncapi/ticks"
)

type ticksSubscriber struct {
	ticks      ticks.Interface
	controller *asyncapi.AppController
}

func newTicksSubscriber(controller *asyncapi.AppController, app ticks.Interface) ticksSubscriber {
	return ticksSubscriber{
		ticks:      app,
		controller: controller,
	}
}

func (s ticksSubscriber) CryptellationTicksRegisterRequest(ctx context.Context, msg asyncapi.RegisteringRequestMessage) {
	log.Printf("Received register request: %+v\n", msg)

	// Set response
	resp := asyncapi.NewRegisteringResponseMessage()
	resp.SetAsResponseFrom(&msg)
	defer func() { _ = s.controller.PublishCryptellationTicksRegisterResponse(ctx, resp) }()

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

func (s ticksSubscriber) CryptellationTicksUnregisterRequest(ctx context.Context, msg asyncapi.RegisteringRequestMessage) {
	log.Printf("Received unregister request: %+v\n", msg)

	// Set response
	resp := asyncapi.NewRegisteringResponseMessage()
	resp.SetAsResponseFrom(&msg)
	defer func() { _ = s.controller.PublishCryptellationTicksUnregisterResponse(ctx, resp) }()

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
