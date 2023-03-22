package nats

import (
	"context"
	"net/http"

	"github.com/digital-feather/cryptellation/internal/exchanges/app"
	"github.com/digital-feather/cryptellation/internal/exchanges/infra/events/nats/generated"
)

type subscriber struct {
	exchanges  app.Controller
	controller *generated.AppController
}

func newSubscriber(controller *generated.AppController, app app.Controller) subscriber {
	return subscriber{
		exchanges:  app,
		controller: controller,
	}
}

func (s subscriber) ExchangesListRequest(msg generated.ExchangesRequestMessage, _ bool) {
	// Prepare response and set send at the end
	resp := generated.NewExchangesResponseMessage()
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
		resp.Payload.Error = &generated.ErrorSchema{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}

	// Add exchanges to response
	resp.Payload.Exchanges = make([]generated.ExchangeSchema, len(exchanges))
	for i, exch := range exchanges {
		// Periods
		periods := make([]generated.PeriodSymbolSchema, len(exch.PeriodsSymbols))
		for j, p := range exch.PeriodsSymbols {
			periods[j] = generated.PeriodSymbolSchema(p)
		}

		// Pairs
		pairs := make([]generated.PairSymbolSchema, len(exch.PairsSymbols))
		for j, p := range exch.PairsSymbols {
			pairs[j] = generated.PairSymbolSchema(p)
		}

		// Exchange
		resp.Payload.Exchanges[i] = generated.ExchangeSchema{
			Fees:         exch.Fees,
			Name:         generated.ExchangeNameSchema(exch.Name),
			Pairs:        pairs,
			Periods:      periods,
			LastSyncTime: exch.LastSyncTime,
		}
	}
}
