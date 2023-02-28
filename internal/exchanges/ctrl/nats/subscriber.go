package nats

import (
	"context"
	"net/http"

	"github.com/digital-feather/cryptellation/internal/exchanges/app"
	"github.com/digital-feather/cryptellation/internal/exchanges/ctrl/nats/internal"
)

type subscriber struct {
	exchanges  app.Controller
	controller *internal.AppController
}

func newSubscriber(controller *internal.AppController, app app.Controller) subscriber {
	return subscriber{
		exchanges:  app,
		controller: controller,
	}
}

func (s subscriber) ExchangesListRequest(msg internal.ExchangesRequestMessage) {
	// Prepare response and set send at the end
	resp := internal.NewExchangesResponseMessage()
	resp.SetAsResponseFrom(msg)
	defer func() { _ = s.controller.PublishExchangesListResponse(resp) }()

	// Change from requests type to application types
	exchangesNames := make([]string, len(msg.Payload))
	for i, e := range msg.Payload {
		exchangesNames[i] = string(e)
	}

	// Request exchange(s) information
	exchanges, err := s.exchanges.GetCached(context.Background(), exchangesNames...)
	if err != nil {
		resp.Payload.Error = &internal.ErrorSchema{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}

	// Add exchanges to response
	resp.Payload.Exchanges = make([]internal.ExchangeSchema, len(exchanges))
	for i, exch := range exchanges {
		// Periods
		periods := make([]internal.PeriodSymbolSchema, len(exch.PeriodsSymbols))
		for j, p := range exch.PeriodsSymbols {
			periods[j] = internal.PeriodSymbolSchema(p)
		}

		// Pairs
		pairs := make([]internal.PairSymbolSchema, len(exch.PairsSymbols))
		for j, p := range exch.PairsSymbols {
			pairs[j] = internal.PairSymbolSchema(p)
		}

		// Exchange
		resp.Payload.Exchanges[i] = internal.ExchangeSchema{
			Fees:         exch.Fees,
			Name:         internal.ExchangeNameSchema(exch.Name),
			Pairs:        pairs,
			Periods:      periods,
			LastSyncTime: exch.LastSyncTime,
		}
	}
}
