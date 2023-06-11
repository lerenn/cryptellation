package events

import (
	"context"
	"net/http"

	"github.com/lerenn/cryptellation/internal/core/exchanges"
)

type subscriber struct {
	exchanges  exchanges.Interface
	controller *AppController
}

func newSubscriber(controller *AppController, app exchanges.Interface) subscriber {
	return subscriber{
		exchanges:  app,
		controller: controller,
	}
}

func (s subscriber) CryptellationExchangesListRequest(msg ExchangesRequestMessage, _ bool) {
	// Prepare response and set send at the end
	resp := NewExchangesResponseMessage()
	resp.SetAsResponseFrom(msg)
	defer func() { _ = s.controller.PublishCryptellationExchangesListResponse(resp) }()

	// Change from requests type to application types
	exchangesNames := msg.ToModel()

	// Request exchange(s) information
	exchanges, err := s.exchanges.GetCached(context.Background(), exchangesNames...)
	if err != nil {
		resp.Payload.Error = &ErrorSchema{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}

	// Add exchanges to response
	resp.Set(exchanges)
}
