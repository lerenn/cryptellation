package nats

import (
	"context"
	"net/http"

	"github.com/lerenn/cryptellation/internal/components/exchanges"
	asyncapi "github.com/lerenn/cryptellation/pkg/asyncapi/exchanges"
)

type exchangesSubscriber struct {
	exchanges  exchanges.Interface
	controller *asyncapi.AppController
}

func newExchangesSubscriber(controller *asyncapi.AppController, app exchanges.Interface) exchangesSubscriber {
	return exchangesSubscriber{
		exchanges:  app,
		controller: controller,
	}
}

func (s exchangesSubscriber) CryptellationExchangesListRequest(ctx context.Context, msg asyncapi.ExchangesRequestMessage) {
	// Prepare response and set send at the end
	resp := asyncapi.NewExchangesResponseMessage()
	resp.SetAsResponseFrom(&msg)
	defer func() { _ = s.controller.PublishCryptellationExchangesListResponse(ctx, resp) }()

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
