// Candlesticks
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.30.2 -g application -p asyncapi -i ../asyncapi.yaml -o ./app.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.30.2 -g user        -p asyncapi -i ../asyncapi.yaml -o ./user.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.30.2 -g types       -p asyncapi -i ../asyncapi.yaml -o ./types.gen.go

package asyncapi

import (
	"time"

	"github.com/lerenn/cryptellation/pkg/utils"
	client "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"
	"github.com/lerenn/cryptellation/svc/candlesticks/internal/app"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"
)

func (msg *ListCandlesticksRequestMessage) Set(payload client.ReadCandlesticksPayload) {
	msg.Payload.ExchangeName = ExchangeNameSchema(payload.ExchangeName)
	msg.Payload.PairSymbol = PairSymbolSchema(payload.PairSymbol)
	msg.Payload.PeriodSymbol = PeriodSymbolSchema(payload.Period.String())
	if payload.Start != nil {
		msg.Payload.Start = utils.ToReference(DateSchema(*payload.Start))
	}
	if payload.End != nil {
		msg.Payload.End = utils.ToReference(DateSchema(*payload.End))
	}
	msg.Payload.Limit = LimitSchema(payload.Limit)
}

func (msg *ListCandlesticksRequestMessage) ToModel() (app.GetCachedPayload, error) {
	// Process specific types
	per, err := period.FromString(string(msg.Payload.PeriodSymbol))
	if err != nil {
		return app.GetCachedPayload{}, err
	}

	// Request list
	return app.GetCachedPayload{
		ExchangeName: string(msg.Payload.ExchangeName),
		PairSymbol:   string(msg.Payload.PairSymbol),
		Period:       per,
		Start:        (*time.Time)(msg.Payload.Start),
		End:          (*time.Time)(msg.Payload.End),
		Limit:        uint(msg.Payload.Limit),
	}, nil
}

func (msg *ListCandlesticksResponseMessage) Set(list *candlestick.List) error {
	respList := make(CandlestickListSchema, 0, list.Len())
	if err := list.Loop(func(t time.Time, cs candlestick.Candlestick) (bool, error) {
		respList = append(respList, CandlestickSchema{
			Time:   DateSchema(t),
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

func (msg *ListCandlesticksResponseMessage) ToModel(exchange, pair string, symbol period.Symbol) (*candlestick.List, error) {
	// Create list
	list := candlestick.NewEmptyList(exchange, pair, symbol)

	// Fill list
	for _, c := range *msg.Payload.Candlesticks {
		if err := list.Set(time.Time(c.Time), candlestick.Candlestick{
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
