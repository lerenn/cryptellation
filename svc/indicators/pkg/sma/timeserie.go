package sma

import (
	"time"

	"github.com/lerenn/cryptellation/pkg/models/timeserie"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
)

type TimeSeriePayload struct {
	Candlesticks *candlestick.List
	PriceType    candlestick.PriceType
	Start        time.Time
	End          time.Time
	PeriodNumber int
}

func TimeSerie(payload TimeSeriePayload) *timeserie.TimeSerie[float64] {
	ts := timeserie.New[float64]()

	// For each theorical point
	for start := payload.Start; payload.End.After(start) || payload.End.Equal(start); start = start.Add(payload.Candlesticks.Period.Duration()) {
		// Get first and last data
		// Note: removing 1 to period number to count the actual time in it
		first := start.Add(-payload.Candlesticks.Period.Duration() * time.Duration(payload.PeriodNumber-1))
		last := start

		// Get interesting candlesticks
		candlesticks := timeserie.New[candlestick.Candlestick]()
		_ = payload.Candlesticks.Loop(func(t time.Time, cs candlestick.Candlestick) (bool, error) {
			if (t.Equal(first) || t.After(first)) && (t.Equal(last) || t.Before(last)) {
				candlesticks.Set(t, cs)
			}

			return false, nil
		})

		// Add calculated point to timeserie
		ts.Set(start, Point(PointPayload{
			Candlesticks: candlesticks,
			PriceType:    payload.PriceType,
		}))
	}

	return ts
}
