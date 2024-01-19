package nats

import (
	"context"
	"net/http"

	"github.com/lerenn/cryptellation/pkg/version"
	asyncapi "github.com/lerenn/cryptellation/svc/exchanges/api/asyncapi"
	exchanges "github.com/lerenn/cryptellation/svc/exchanges/internal/app"
)

type subscriber struct {
	exchanges  exchanges.Exchanges
	controller *asyncapi.AppController
}

func newSubscriber(controller *asyncapi.AppController, app exchanges.Exchanges) subscriber {
	return subscriber{
		exchanges:  app,
		controller: controller,
	}
}

func (s subscriber) ListExchangesRequest(ctx context.Context, msg asyncapi.ListExchangesRequestMessage) {
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

func (s subscriber) ServiceInfoRequest(ctx context.Context, msg asyncapi.ServiceInfoRequestMessage) {
	// Prepare response and set send at the end
	resp := asyncapi.NewServiceInfoResponseMessage()
	resp.SetAsResponseFrom(&msg)
	defer func() { _ = s.controller.PublishServiceInfoResponse(ctx, resp) }()

	// Set info
	resp.Payload.ApiVersion = asyncapi.AsyncAPIVersion
	resp.Payload.BinVersion = version.Version()
}
