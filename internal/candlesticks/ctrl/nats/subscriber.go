package nats

import (
	"context"
	"net/http"

	asyncapi "github.com/digital-feather/cryptellation/api/asyncapi/candlesticks"
	"github.com/digital-feather/cryptellation/internal/candlesticks/app"
)

type subscriber struct {
	candlesticks app.Controller
	controller   *asyncapi.AppController
}

func newSubscriber(controller *asyncapi.AppController, app app.Controller) subscriber {
	return subscriber{
		candlesticks: app,
		controller:   controller,
	}
}

func (s subscriber) CandlesticksListRequest(msg asyncapi.CandlesticksListRequestMessage, _ bool) {
	// Prepare response and set send at the end
	resp := asyncapi.NewCandlesticksListResponseMessage()
	resp.SetAsResponseFrom(msg)
	defer func() { _ = s.controller.PublishCandlesticksListResponse(resp) }()

	// Set list payload
	payload, err := msg.ToModel()
	if err != nil {
		resp.Payload.Error = &asyncapi.ErrorSchema{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		return
	}

	// Request list
	list, err := s.candlesticks.GetCached(context.Background(), payload)
	if err != nil {
		resp.Payload.Error = &asyncapi.ErrorSchema{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}

	// Add list to response
	if err := resp.Set(list); err != nil {
		resp.Payload.Error = &asyncapi.ErrorSchema{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}
}
