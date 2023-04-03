package candlesticks

import (
	"context"
	"net/http"

	"github.com/digital-feather/cryptellation/services/candlesticks"
)

type subscriber struct {
	candlesticks candlesticks.Interface
	controller   *AppController
}

func newSubscriber(controller *AppController, app candlesticks.Interface) subscriber {
	return subscriber{
		candlesticks: app,
		controller:   controller,
	}
}

func (s subscriber) CryptellationCandlesticksListRequest(msg CandlesticksListRequestMessage, _ bool) {
	// Prepare response and set send at the end
	resp := NewCandlesticksListResponseMessage()
	resp.SetAsResponseFrom(msg)
	defer func() { _ = s.controller.PublishCryptellationCandlesticksListResponse(resp) }()

	// Set list payload
	payload, err := msg.ToModel()
	if err != nil {
		resp.Payload.Error = &ErrorSchema{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		return
	}

	// Request list
	list, err := s.candlesticks.GetCached(context.Background(), payload)
	if err != nil {
		resp.Payload.Error = &ErrorSchema{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}

	// Add list to response
	if err := resp.Set(list); err != nil {
		resp.Payload.Error = &ErrorSchema{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}
}
