package nats

import (
	"context"
	"net/http"

	asyncapi "github.com/digital-feather/cryptellation/api/asyncapi/exchanges"
	"github.com/digital-feather/cryptellation/internal/exchanges/app"
)

type subscriber struct {
	exchanges  app.Controller
	controller *asyncapi.AppController
}

func newSubscriber(controller *asyncapi.AppController, app app.Controller) subscriber {
	return subscriber{
		exchanges:  app,
		controller: controller,
	}
}

func (s subscriber) ExchangesListRequest(msg asyncapi.ExchangesRequestMessage, _ bool) {
	// Prepare response and set send at the end
	resp := asyncapi.NewExchangesResponseMessage()
	resp.SetAsResponseFrom(msg)
	defer func() { _ = s.controller.PublishExchangesListResponse(resp) }()

	// Change from requests type to application types
	exchangesNames := msg.ToModel()

	// Request exchange(s) information
	exchanges, err := s.exchanges.GetCached(context.Background(), exchangesNames...)
	if err != nil {
		resp.Payload.Error = &asyncapi.ErrorSchema{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}

	// Add exchanges to response
	resp.Set(exchanges)
}
