package ticks

import (
	"context"
	"log"
	"net/http"

	"github.com/lerenn/cryptellation/services/ticks"
)

type subscriber struct {
	ticks      ticks.Interface
	controller *AppController
}

func newSubscriber(controller *AppController, app ticks.Interface) subscriber {
	return subscriber{
		ticks:      app,
		controller: controller,
	}
}

func (s subscriber) CryptellationTicksRegisterRequest(msg RegisteringRequestMessage, _ bool) {
	log.Printf("Received register request: %+v\n", msg)

	// Set response
	resp := NewRegisteringResponseMessage()
	resp.SetAsResponseFrom(msg)
	defer func() { _ = s.controller.PublishCryptellationTicksRegisterResponse(resp) }()

	// Register as requested
	count, err := s.ticks.Register(
		context.Background(),
		string(msg.Payload.Exchange),
		string(msg.Payload.Pair),
	)

	// If there is an error, return it as BadRequest
	// TODO: get correct error
	if err != nil {
		resp.Payload.Error = &ErrorSchema{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		return
	}

	// Otherwise, return count
	resp.Payload.Count = &count
}

func (s subscriber) CryptellationTicksUnregisterRequest(msg RegisteringRequestMessage, _ bool) {
	log.Printf("Received unregister request: %+v\n", msg)

	// Set response
	resp := NewRegisteringResponseMessage()
	resp.SetAsResponseFrom(msg)
	defer func() { _ = s.controller.PublishCryptellationTicksUnregisterResponse(resp) }()

	// Register as requested
	count, err := s.ticks.Unregister(
		context.Background(),
		string(msg.Payload.Exchange),
		string(msg.Payload.Pair),
	)

	// If there is an error, return it as BadRequest
	// TODO: get correct error
	if err != nil {
		resp.Payload.Error = &ErrorSchema{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		return
	}

	// Otherwise, return count
	resp.Payload.Count = &count
}
