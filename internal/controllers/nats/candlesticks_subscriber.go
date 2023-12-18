package nats

import (
	"context"
	"net/http"

	asyncapi "github.com/lerenn/cryptellation/api/asyncapi/candlesticks"
	"github.com/lerenn/cryptellation/internal/components/candlesticks"
	"github.com/lerenn/cryptellation/pkg/version"
)

type candlesticksSubscriber struct {
	candlesticks candlesticks.Interface
	controller   *asyncapi.AppController
}

func newCandlesticksSubscriber(controller *asyncapi.AppController, app candlesticks.Interface) candlesticksSubscriber {
	return candlesticksSubscriber{
		candlesticks: app,
		controller:   controller,
	}
}

func (s candlesticksSubscriber) ListCandlesticksRequest(ctx context.Context, msg asyncapi.ListCandlesticksRequestMessage) {
	// Prepare response and set send at the end
	resp := asyncapi.NewListCandlesticksResponseMessage()
	resp.SetAsResponseFrom(&msg)
	defer func() { _ = s.controller.PublishListCandlesticksResponse(ctx, resp) }()

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

func (s candlesticksSubscriber) ServiceInfoRequest(ctx context.Context, msg asyncapi.ServiceInfoRequestMessage) {
	// Prepare response and set send at the end
	resp := asyncapi.NewServiceInfoResponseMessage()
	resp.SetAsResponseFrom(&msg)
	defer func() { _ = s.controller.PublishServiceInfoResponse(ctx, resp) }()

	// Set info
	resp.Payload.ApiVersion = asyncapi.AsyncAPIVersion
	resp.Payload.BinVersion = version.GetVersion()
}
