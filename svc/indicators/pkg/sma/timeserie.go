package sma

import (
	"time"

	"github.com/lerenn/cryptellation/pkg/models/timeserie"

	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
)

type TimeSeriePayload struct {
	Candlesticks *candlestick.List
	PriceType    candlestick.Price
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
		_ = payload.Candlesticks.Loop(func(cs candlestick.Candlestick) (bool, error) {
			if (cs.Time.Equal(first) || cs.Time.After(first)) && (cs.Time.Equal(last) || cs.Time.Before(last)) {
				candlesticks.Set(cs.Time, cs)
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

func InvalidValues(ts *timeserie.TimeSerie[float64]) bool {
	invalidValuesDetected := false
	_ = ts.Loop(func(t time.Time, v float64) (bool, error) {
		if v == 0 {
			invalidValuesDetected = true
			return true, nil
		}
		return false, nil
	})
	return invalidValuesDetected
}
