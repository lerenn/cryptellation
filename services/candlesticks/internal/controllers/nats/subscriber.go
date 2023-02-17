package nats

import (
	"context"
	"net/http"
	"time"

	"github.com/digital-feather/cryptellation/pkg/types/candlestick"
	"github.com/digital-feather/cryptellation/pkg/types/period"
	async "github.com/digital-feather/cryptellation/services/candlesticks/internal/adapters/async"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/components/candlesticks"
)

type subscriber struct {
	candlesticks candlesticks.Port
	controller   *async.AppController
}

func newSubscriber(controller *async.AppController, app candlesticks.Port) subscriber {
	return subscriber{
		candlesticks: app,
		controller:   controller,
	}
}

func (s subscriber) CandlesticksListRequest(msg async.CandlesticksListRequestMessage) {
	// Prepare response and set send at the end
	resp := async.NewCandlesticksListResponseMessage()
	resp.SetAsResponseFrom(msg)
	defer func() { _ = s.controller.PublishCandlesticksListResponse(resp) }()

	// Process specific types
	per, err := period.FromString(string(msg.Payload.PeriodSymbol))
	if err != nil {
		resp.Payload.Error = &async.ErrorSchema{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		return
	}

	// Request list
	list, err := s.candlesticks.GetCached(context.Background(), candlesticks.GetCachedPayload{
		ExchangeName: string(msg.Payload.ExchangeName),
		PairSymbol:   string(msg.Payload.PairSymbol),
		Period:       per,
		Start:        (*time.Time)(msg.Payload.Start),
		End:          (*time.Time)(msg.Payload.End),
		Limit:        uint(msg.Payload.Limit),
	})
	if err != nil {
		resp.Payload.Error = &async.ErrorSchema{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}

	// Add list to response
	respList := make(async.CandlestickListSchema, 0, list.Len())
	if err := list.Loop(func(t time.Time, cs candlestick.Candlestick) (bool, error) {
		respList = append(respList, async.CandlestickSchema{
			Time:   async.DateSchema(t),
			Open:   &cs.Open,
			High:   &cs.High,
			Low:    &cs.Low,
			Close:  &cs.Close,
			Volume: &cs.Volume,
		})
		return false, nil
	}); err != nil {
		resp.Payload.Error = &async.ErrorSchema{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}
	resp.Payload.Candlesticks = &respList
}
