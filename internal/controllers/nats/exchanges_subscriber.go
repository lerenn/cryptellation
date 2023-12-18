package nats

import (
	"context"
	"net/http"

	asyncapi "github.com/lerenn/cryptellation/api/asyncapi/exchanges"
	"github.com/lerenn/cryptellation/internal/components/exchanges"
	"github.com/lerenn/cryptellation/pkg/version"
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

func (s exchangesSubscriber) ListExchangesRequest(ctx context.Context, msg asyncapi.ListExchangesRequestMessage) {
	// Prepare response and set send at the end
	resp := asyncapi.NewListExchangesResponseMessage()
	resp.SetAsResponseFrom(&msg)
	defer func() { _ = s.controller.PublishListExchangesResponse(ctx, resp) }()

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

func (s exchangesSubscriber) ServiceInfoRequest(ctx context.Context, msg asyncapi.ServiceInfoRequestMessage) {
	// Prepare response and set send at the end
	resp := asyncapi.NewServiceInfoResponseMessage()
	resp.SetAsResponseFrom(&msg)
	defer func() { _ = s.controller.PublishServiceInfoResponse(ctx, resp) }()

	// Set info
	resp.Payload.ApiVersion = asyncapi.AsyncAPIVersion
	resp.Payload.BinVersion = version.GetVersion()
}
