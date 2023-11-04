package nats

import (
	"context"
	"net/http"

	"github.com/lerenn/cryptellation/internal/components/candlesticks"
	asyncapi "github.com/lerenn/cryptellation/pkg/asyncapi/candlesticks"
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

func (s candlesticksSubscriber) CryptellationCandlesticksListRequest(ctx context.Context, msg asyncapi.CandlesticksListRequestMessage) {
	// Prepare response and set send at the end
	resp := asyncapi.NewCandlesticksListResponseMessage()
	resp.SetAsResponseFrom(&msg)
	defer func() { _ = s.controller.PublishCryptellationCandlesticksListResponse(ctx, resp) }()

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
