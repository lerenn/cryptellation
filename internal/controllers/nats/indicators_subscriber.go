package nats

import (
	"context"
	"fmt"
	"net/http"

	asyncapi "github.com/lerenn/cryptellation/api/asyncapi/indicators"
	"github.com/lerenn/cryptellation/internal/components/indicators"
	"github.com/lerenn/cryptellation/pkg/version"
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

func (s indicatorsSubscriber) GetSMARequest(ctx context.Context, msg asyncapi.GetSMARequestMessage) {
	// Prepare response and set send at the end
	resp := asyncapi.NewGetSMAResponseMessage()
	resp.SetAsResponseFrom(&msg)
	defer func() { _ = s.controller.PublishGetSMAResponse(ctx, resp) }()

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

func (s indicatorsSubscriber) ServiceInfoRequest(ctx context.Context, msg asyncapi.ServiceInfoRequestMessage) {
	// Prepare response and set send at the end
	resp := asyncapi.NewServiceInfoResponseMessage()
	resp.SetAsResponseFrom(&msg)
	defer func() { _ = s.controller.PublishServiceInfoResponse(ctx, resp) }()

	// Set info
	resp.Payload.ApiVersion = asyncapi.AsyncAPIVersion
	resp.Payload.BinVersion = version.GetVersion()
}
