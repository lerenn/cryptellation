package events

import (
	"context"
	"fmt"
	"net/http"

	"github.com/lerenn/cryptellation/internal/core/indicators"
)

type subscriber struct {
	indicators indicators.Interface
	controller *AppController
}

func newSubscriber(controller *AppController, app indicators.Interface) subscriber {
	return subscriber{
		indicators: app,
		controller: controller,
	}
}

func (s subscriber) CryptellationIndicatorsSmaRequest(msg SmaRequestMessage, _ bool) {
	// Prepare response and set send at the end
	resp := NewSmaResponseMessage()
	resp.SetAsResponseFrom(msg)
	defer func() { _ = s.controller.PublishCryptellationIndicatorsSmaResponse(resp) }()

	// Change from requests type to application types
	payload, err := msg.ToModel()
	if err != nil {
		resp.Payload.Error = &ErrorSchema{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		return
	}

	// Request exchange(s) information
	indicators, err := s.indicators.GetCachedSMA(context.Background(), payload)
	if err != nil {
		resp.Payload.Error = &ErrorSchema{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}

	// Add indicators to response
	resp.Set(indicators)
	fmt.Println(len(*resp.Payload.Data))
}
