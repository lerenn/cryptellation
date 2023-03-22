package nats

import (
	"context"
	"log"
	"net/http"

	asyncapi "github.com/digital-feather/cryptellation/api/asyncapi/ticks"
	"github.com/digital-feather/cryptellation/internal/ticks/app"
)

type subscriber struct {
	ticks      app.Controller
	controller *asyncapi.AppController
}

func newSubscriber(controller *asyncapi.AppController, app app.Controller) subscriber {
	return subscriber{
		ticks:      app,
		controller: controller,
	}
}

func (s subscriber) TicksRegisterRequest(msg asyncapi.RegisteringRequestMessage, _ bool) {
	log.Printf("Received register request: %+v\n", msg)

	// Set response
	resp := asyncapi.NewRegisteringResponseMessage()
	resp.SetAsResponseFrom(msg)
	defer func() { _ = s.controller.PublishTicksRegisterResponse(resp) }()

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

func (s subscriber) TicksUnregisterRequest(msg asyncapi.RegisteringRequestMessage, _ bool) {
	log.Printf("Received unregister request: %+v\n", msg)

	// Set response
	resp := asyncapi.NewRegisteringResponseMessage()
	resp.SetAsResponseFrom(msg)
	defer func() { _ = s.controller.PublishTicksUnregisterResponse(resp) }()

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
