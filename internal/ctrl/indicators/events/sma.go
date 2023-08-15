package events

import (
	"time"

	client "github.com/lerenn/cryptellation/clients/go"
	"github.com/lerenn/cryptellation/internal/core/indicators"
	"github.com/lerenn/cryptellation/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/pkg/models/period"
	"github.com/lerenn/cryptellation/pkg/models/timeserie"
	"github.com/lerenn/cryptellation/pkg/utils"
)

func (msg *SmaRequestMessage) Set(payload client.SMAPayload) {
	msg.Payload.ExchangeName = ExchangeNameSchema(payload.ExchangeName)
	msg.Payload.PairSymbol = PairSymbolSchema(payload.PairSymbol)
	msg.Payload.PeriodSymbol = PeriodSymbolSchema(payload.Period)
	msg.Payload.Start = utils.ToReference(DateSchema(payload.Start))
	msg.Payload.End = utils.ToReference(DateSchema(payload.End))
	msg.Payload.PeriodNumber = NumberOfPeriodsSchema(payload.PeriodNumber)
	msg.Payload.PriceType = utils.ToReference(PriceTypeSchema(payload.PriceType))
}

func (msg *SmaRequestMessage) ToModel() (indicators.GetCachedSMAPayload, error) {
	per := period.Symbol(msg.Payload.PeriodSymbol)
	if err := per.Validate(); err != nil {
		return indicators.GetCachedSMAPayload{}, err
	}

	pt := candlestick.PriceType(*msg.Payload.PriceType)
	if err := pt.Validate(); err != nil {
		return indicators.GetCachedSMAPayload{}, err
	}

	return indicators.GetCachedSMAPayload{
		ExchangeName: string(msg.Payload.ExchangeName),
		PairSymbol:   string(msg.Payload.PairSymbol),
		Period:       per,
		Start:        time.Time(*msg.Payload.Start),
		End:          time.Time(*msg.Payload.End),
		PeriodNumber: uint(msg.Payload.PeriodNumber),
		PriceType:    pt,
	}, nil
}

func (msg *SmaResponseMessage) Set(ts *timeserie.TimeSerie[float64]) {
	count := 0
	data := make(NumericTimeSerieSchema, ts.Len())
	_ = ts.Loop(func(t time.Time, v float64) (bool, error) {
		point := data[count]

		// Set point
		point.Time = DateSchema(t)
		point.Value = v
		data[count] = point

		count++
		return false, nil
	})

	msg.Payload.Data = &data
}

func (msg *SmaResponseMessage) ToModel() *timeserie.TimeSerie[float64] {
	ts := timeserie.New[float64]()
	for _, point := range *msg.Payload.Data {
		ts.Set(time.Time(point.Time), point.Value)
	}
	return ts
}
