package nats

import (
	"context"
	"net/http"
	"time"

	"github.com/digital-feather/cryptellation/internal/candlesticks/app"
	"github.com/digital-feather/cryptellation/internal/candlesticks/ctrl/nats/internal"
	"github.com/digital-feather/cryptellation/pkg/candlestick"
	"github.com/digital-feather/cryptellation/pkg/period"
)

type subscriber struct {
	candlesticks app.Port
	controller   *internal.AppController
}

func newSubscriber(controller *internal.AppController, app app.Port) subscriber {
	return subscriber{
		candlesticks: app,
		controller:   controller,
	}
}

func (s subscriber) CandlesticksListRequest(msg internal.CandlesticksListRequestMessage) {
	// Prepare response and set send at the end
	resp := internal.NewCandlesticksListResponseMessage()
	resp.SetAsResponseFrom(msg)
	defer func() { _ = s.controller.PublishCandlesticksListResponse(resp) }()

	// Process specific types
	per, err := period.FromString(string(msg.Payload.PeriodSymbol))
	if err != nil {
		resp.Payload.Error = &internal.ErrorSchema{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		return
	}

	// Request list
	list, err := s.candlesticks.GetCached(context.Background(), app.GetCachedPayload{
		ExchangeName: string(msg.Payload.ExchangeName),
		PairSymbol:   string(msg.Payload.PairSymbol),
		Period:       per,
		Start:        (*time.Time)(msg.Payload.Start),
		End:          (*time.Time)(msg.Payload.End),
		Limit:        uint(msg.Payload.Limit),
	})
	if err != nil {
		resp.Payload.Error = &internal.ErrorSchema{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}

	// Add list to response
	respList := make(internal.CandlestickListSchema, 0, list.Len())
	if err := list.Loop(func(t time.Time, cs candlestick.Candlestick) (bool, error) {
		respList = append(respList, internal.CandlestickSchema{
			Time:   internal.DateSchema(t),
			Open:   cs.Open,
			High:   cs.High,
			Low:    cs.Low,
			Close:  cs.Close,
			Volume: cs.Volume,
		})
		return false, nil
	}); err != nil {
		resp.Payload.Error = &internal.ErrorSchema{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}
	resp.Payload.Candlesticks = &respList
}
