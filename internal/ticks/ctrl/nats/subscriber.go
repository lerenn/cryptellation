package nats

import (
	"context"
	"log"
	"net/http"

	"github.com/digital-feather/cryptellation/internal/ticks/app"
	"github.com/digital-feather/cryptellation/internal/ticks/infra/events/nats/generated"
)

type subscriber struct {
	ticks      app.Controller
	controller *generated.AppController
}

func newSubscriber(controller *generated.AppController, app app.Controller) subscriber {
	return subscriber{
		ticks:      app,
		controller: controller,
	}
}

func (s subscriber) TicksRegisterRequest(msg generated.RegisteringRequestMessage, _ bool) {
	log.Printf("Received register request: %+v\n", msg)

	// Set response
	resp := generated.NewRegisteringResponseMessage()
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
		resp.Payload.Error = &generated.ErrorSchema{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		return
	}

	// Otherwise, return count
	resp.Payload.Count = &count
}

func (s subscriber) TicksUnregisterRequest(msg generated.RegisteringRequestMessage, _ bool) {
	log.Printf("Received unregister request: %+v\n", msg)

	// Set response
	resp := generated.NewRegisteringResponseMessage()
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
		resp.Payload.Error = &generated.ErrorSchema{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		return
	}

	// Otherwise, return count
	resp.Payload.Count = &count
}
