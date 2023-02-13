package nats

import (
	"context"
	"net/http"
	"time"

	"github.com/digital-feather/cryptellation/services/candlesticks/internal/application"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/application/operators/candlesticks"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/controllers/nats/gen"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/candlestick"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/period"
)

type subscriber struct {
	app        *application.Application
	controller *gen.AppController
}

func newSubscriber(controller *gen.AppController, app *application.Application) subscriber {
	return subscriber{
		app:        app,
		controller: controller,
	}
}

func (s subscriber) CandlesticksListRequest(msg gen.CandlesticksListRequestMessage) {
	// Prepare response and set send at the end
	resp := gen.NewCandlesticksListResponseMessage()
	resp.SetAsResponseFrom(msg)
	defer s.controller.PublishCandlesticksListResponse(resp)

	// Process specific types
	per, err := period.FromString(string(msg.Payload.PeriodSymbol))
	if err != nil {
		resp.Payload.Error = &gen.ErrorSchema{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		return
	}

	// Request list
	list, err := s.app.Candlesticks.GetCached(context.Background(), candlesticks.GetCachedPayload{
		ExchangeName: string(msg.Payload.ExchangeName),
		PairSymbol:   string(msg.Payload.PairSymbol),
		Period:       per,
		Start:        (*time.Time)(msg.Payload.Start),
		End:          (*time.Time)(msg.Payload.End),
		Limit:        uint(msg.Payload.Limit),
	})
	if err != nil {
		resp.Payload.Error = &gen.ErrorSchema{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}

	// Add list to response
	respList := make(gen.CandlestickListSchema, 0, list.Len())
	if err := list.Loop(func(t time.Time, cs candlestick.Candlestick) (bool, error) {
		respList = append(respList, gen.CandlestickSchema{
			Time:   gen.DateSchema(t),
			Open:   &cs.Open,
			High:   &cs.High,
			Low:    &cs.Low,
			Close:  &cs.Close,
			Volume: &cs.Volume,
		})
		return false, nil
	}); err != nil {
		resp.Payload.Error = &gen.ErrorSchema{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}
	resp.Payload.Candlesticks = &respList
}
