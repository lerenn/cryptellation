package nats

import (
	"context"
	"fmt"
	"net/http"

	"github.com/lerenn/cryptellation/pkg/version"
	asyncapi "github.com/lerenn/cryptellation/svc/indicators/api/asyncapi"
	"github.com/lerenn/cryptellation/svc/indicators/internal/app"
)

type subscriber struct {
	indicators app.Indicators
	controller *asyncapi.AppController
}

func newSubscriber(controller *asyncapi.AppController, app app.Indicators) subscriber {
	return subscriber{
		indicators: app,
		controller: controller,
	}
}

func (s subscriber) GetSMARequest(ctx context.Context, msg asyncapi.GetSMARequestMessage) {
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

func (s subscriber) ServiceInfoRequest(ctx context.Context, msg asyncapi.ServiceInfoRequestMessage) {
	// Prepare response and set send at the end
	resp := asyncapi.NewServiceInfoResponseMessage()
	resp.SetAsResponseFrom(&msg)
	defer func() { _ = s.controller.PublishServiceInfoResponse(ctx, resp) }()

	// Set info
	resp.Payload.ApiVersion = asyncapi.AsyncAPIVersion
	resp.Payload.BinVersion = version.Version()
}
