// Candlesticks
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.13.1 -g application -p candlesticks -i ./../candlesticks.yaml -o ./app.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.13.1 -g client      -p candlesticks -i ./../candlesticks.yaml -o ./client.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.13.1 -g broker      -p candlesticks -i ./../candlesticks.yaml -o ./broker.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.13.1 -g types       -p candlesticks -i ./../candlesticks.yaml -o ./types.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.13.1 -g nats        -p candlesticks -i ./../candlesticks.yaml -o ./nats.gen.go

package candlesticks

import (
	"time"

	client "github.com/digital-feather/cryptellation/clients/go"
	"github.com/digital-feather/cryptellation/pkg/candlestick"
	"github.com/digital-feather/cryptellation/pkg/period"
	"github.com/digital-feather/cryptellation/pkg/utils"
	"github.com/digital-feather/cryptellation/services/candlesticks"
)

func (msg *CandlesticksListRequestMessage) Set(payload client.ReadCandlesticksPayload) {
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

func (msg *CandlesticksListRequestMessage) ToModel() (candlesticks.GetCachedPayload, error) {
	// Process specific types
	per, err := period.FromString(string(msg.Payload.PeriodSymbol))
	if err != nil {
		return candlesticks.GetCachedPayload{}, err
	}

	// Request list
	return candlesticks.GetCachedPayload{
		ExchangeName: string(msg.Payload.ExchangeName),
		PairSymbol:   string(msg.Payload.PairSymbol),
		Period:       per,
		Start:        (*time.Time)(msg.Payload.Start),
		End:          (*time.Time)(msg.Payload.End),
		Limit:        uint(msg.Payload.Limit),
	}, nil
}

func (msg *CandlesticksListResponseMessage) Set(list *candlestick.List) error {
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

func (msg *CandlesticksListResponseMessage) ToModel(exchange, pair string, symbol period.Symbol) (*candlestick.List, error) {
	// Create list
	list := candlestick.NewEmptyList(candlestick.ListID{
		ExchangeName: exchange,
		PairSymbol:   pair,
		Period:       symbol,
	})

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
