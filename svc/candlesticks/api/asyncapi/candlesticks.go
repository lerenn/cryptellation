// Candlesticks
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.39.0 -g application -p asyncapi -i ../asyncapi.yaml,../../../../internal/asyncapi/models.yaml -o ./app.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.39.0 -g user        -p asyncapi -i ../asyncapi.yaml,../../../../internal/asyncapi/models.yaml -o ./user.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.39.0 -g types       -p asyncapi -i ../asyncapi.yaml,../../../../internal/asyncapi/models.yaml -o ./types.gen.go

package asyncapi

import (
	"time"

	"github.com/lerenn/cryptellation/pkg/utils"

	client "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"
	"github.com/lerenn/cryptellation/svc/candlesticks/internal/app"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"
)

func (msg *ListRequestMessage) Set(payload client.ReadCandlesticksPayload) {
	msg.Payload.Exchange = ExchangeSchema(payload.Exchange)
	msg.Payload.Pair = PairSchema(payload.Pair)
	msg.Payload.Period = PeriodSchema(payload.Period.String())
	if payload.Start != nil {
		msg.Payload.Start = utils.ToReference(DateSchema(*payload.Start))
	}
	if payload.End != nil {
		msg.Payload.End = utils.ToReference(DateSchema(*payload.End))
	}
	msg.Payload.Limit = LimitSchema(payload.Limit)
}

func (msg *ListRequestMessage) ToModel() (app.GetCachedPayload, error) {
	// Process specific types
	per, err := period.FromString(string(msg.Payload.Period))
	if err != nil {
		return app.GetCachedPayload{}, err
	}

	// Request list
	return app.GetCachedPayload{
		Exchange: string(msg.Payload.Exchange),
		Pair:     string(msg.Payload.Pair),
		Period:   per,
		Start:    (*time.Time)(msg.Payload.Start),
		End:      (*time.Time)(msg.Payload.End),
		Limit:    uint(msg.Payload.Limit),
	}, nil
}

func (msg *ListResponseMessage) Set(list *candlestick.List) error {
	respList := make(CandlestickListSchema, 0, list.Len())
	if err := list.Loop(func(cs candlestick.Candlestick) (bool, error) {
		respList = append(respList, CandlestickSchema{
			Time:   DateSchema(cs.Time),
			Open:   cs.Open,
			High:   cs.High,
			Low:    cs.Low,
			Close:  cs.Close,
			Volume: cs.Volume,
		})
		return false, nil
	}); err != nil {
		return err
	}

	msg.Payload.Candlesticks = &respList
	return nil
}

func (msg *ListResponseMessage) ToModel(exchange, pair string, symbol period.Symbol) (*candlestick.List, error) {
	// Create list
	list := candlestick.NewList(exchange, pair, symbol)

	// Fill list
	for _, c := range *msg.Payload.Candlesticks {
		if err := list.Set(candlestick.Candlestick{
			Time:   time.Time(c.Time),
			Open:   c.Open,
			High:   c.High,
			Low:    c.Low,
			Close:  c.Close,
			Volume: c.Volume,
		}); err != nil {
			return nil, err
		}
	}

	return list, nil
}
