package app

import (
	"context"
	"time"

	"cryptellation/pkg/models/timeserie"

	"cryptellation/svc/candlesticks/pkg/candlestick"
	"cryptellation/svc/candlesticks/pkg/period"
)

type Indicators interface {
	GetCachedSMA(ctx context.Context, payload GetCachedSMAPayload) (*timeserie.TimeSerie[float64], error)
}

type GetCachedSMAPayload struct {
	Exchange     string
	Pair         string
	Period       period.Symbol
	Start        time.Time
	End          time.Time
	PeriodNumber int
	PriceType    candlestick.PriceType
}

func (payload *GetCachedSMAPayload) Process() {
	// Round time
	payload.Start = payload.Period.RoundTime(payload.Start)
	payload.End = payload.Period.RoundTime(payload.End)
}
