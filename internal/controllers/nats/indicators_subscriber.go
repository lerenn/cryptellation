package nats

import (
	"context"
	"fmt"
	"net/http"

	"github.com/lerenn/cryptellation/internal/components/indicators"
	asyncapi "github.com/lerenn/cryptellation/pkg/asyncapi/indicators"
)

type indicatorsSubscriber struct {
	indicators indicators.Interface
	controller *asyncapi.AppController
}

func newIndicatorsSubscriber(controller *asyncapi.AppController, app indicators.Interface) indicatorsSubscriber {
	return indicatorsSubscriber{
		indicators: app,
		controller: controller,
	}
}

func (s indicatorsSubscriber) CryptellationIndicatorsSmaRequest(ctx context.Context, msg asyncapi.SmaRequestMessage) {
	// Prepare response and set send at the end
	resp := asyncapi.NewSmaResponseMessage()
	resp.SetAsResponseFrom(&msg)
	defer func() { _ = s.controller.PublishCryptellationIndicatorsSmaResponse(ctx, resp) }()

	// Change from requests type to application types
	payload, err := msg.ToModel()
	if err != nil {
		resp.Payload.Error = &asyncapi.ErrorSchema{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		return
	}

	// Request exchange(s) information
	indicators, err := s.indicators.GetCachedSMA(context.Background(), payload)
	if err != nil {
		resp.Payload.Error = &asyncapi.ErrorSchema{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}

	// Add indicators to response
	resp.Set(indicators)
	fmt.Println(len(*resp.Payload.Data))
}
