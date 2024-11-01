package app

import (
	"context"
	"time"

	"github.com/lerenn/cryptellation/pkg/models/timeserie"

	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"
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
	PriceType    candlestick.Price
}

func (payload *GetCachedSMAPayload) Process() {
	// Round time
	payload.Start = payload.Period.RoundTime(payload.Start)
	payload.End = payload.Period.RoundTime(payload.End)
}
