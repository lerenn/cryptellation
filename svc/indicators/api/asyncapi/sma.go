package asyncapi

import (
	"math"
	"time"

	"github.com/lerenn/cryptellation/pkg/timeserie"
	"github.com/lerenn/cryptellation/pkg/utils"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"
	client "github.com/lerenn/cryptellation/svc/indicators/clients/go"
	"github.com/lerenn/cryptellation/svc/indicators/internal/app"
)

func (msg *SMARequestMessage) Set(payload client.SMAPayload) {
	msg.Payload.Exchange = ExchangeSchema(payload.Exchange)
	msg.Payload.Pair = PairSchema(payload.Pair)
	msg.Payload.Period = PeriodSchema(payload.Period)
	msg.Payload.Start = utils.ToReference(DateSchema(payload.Start))
	msg.Payload.End = utils.ToReference(DateSchema(payload.End))
	msg.Payload.PeriodNumber = NumberOfPeriodsSchema(payload.PeriodNumber)
	msg.Payload.PriceType = utils.ToReference(PriceTypeSchema(payload.PriceType))
}

func (msg *SMARequestMessage) ToModel() (app.GetCachedSMAPayload, error) {
	per := period.Symbol(msg.Payload.Period)
	if err := per.Validate(); err != nil {
		return app.GetCachedSMAPayload{}, err
	}

	pt := candlestick.PriceType(*msg.Payload.PriceType)
	if err := pt.Validate(); err != nil {
		return app.GetCachedSMAPayload{}, err
	}

	return app.GetCachedSMAPayload{
		Exchange:     string(msg.Payload.Exchange),
		Pair:         string(msg.Payload.Pair),
		Period:       per,
		Start:        time.Time(*msg.Payload.Start),
		End:          time.Time(*msg.Payload.End),
		PeriodNumber: int(msg.Payload.PeriodNumber),
		PriceType:    pt,
	}, nil
}

func (msg *SMAResponseMessage) Set(ts *timeserie.TimeSerie[float64]) {
	count := 0
	data := make(NumericTimeSerieSchema, ts.Len())
	_ = ts.Loop(func(t time.Time, v float64) (bool, error) {
		point := data[count]

		// Set point
		point.Time = DateSchema(t)
		if !math.IsNaN(v) {
			point.Value = &v
		}
		data[count] = point

		count++
		return false, nil
	})

	msg.Payload.Data = &data
}

func (msg *SMAResponseMessage) ToModel() *timeserie.TimeSerie[float64] {
	ts := timeserie.New[float64]()
	for _, point := range *msg.Payload.Data {
		if point.Value != nil {
			ts.Set(time.Time(point.Time), *point.Value)
		}
	}
	return ts
}
