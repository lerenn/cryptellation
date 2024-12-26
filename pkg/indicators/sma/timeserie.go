package sma

import (
	"time"

	"github.com/lerenn/cryptellation/v1/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/v1/pkg/models/timeserie"
)

// TimeSerieParams is the parameters needed to create a timeserie.
type TimeSerieParams struct {
	Candlesticks *candlestick.List
	PriceType    candlestick.PriceType
	Start        time.Time
	End          time.Time
	PeriodNumber int
}

// TimeSerie returns a timeserie of calculated points.
func TimeSerie(params TimeSerieParams) (*timeserie.TimeSerie[float64], error) {
	ts := timeserie.New[float64]()

	// For each theorical point
	duration := params.Candlesticks.Metadata.Period.Duration()
	for start := params.Start; params.End.After(start) || params.End.Equal(start); start = start.Add(duration) {
		// Get first and last data
		// Note: removing 1 to period number to count the actual time in it
		first := start.Add(-duration * time.Duration(params.PeriodNumber-1))
		last := start

		// Get interesting candlesticks
		candlesticks := candlestick.NewListWithMetadata(params.Candlesticks.Metadata)
		if err := params.Candlesticks.Loop(func(cs candlestick.Candlestick) (bool, error) {
			if (cs.Time.Equal(first) || cs.Time.After(first)) && (cs.Time.Equal(last) || cs.Time.Before(last)) {
				if err := candlesticks.Set(cs); err != nil {
					return true, err
				}
			}

			return false, nil
		}); err != nil {
			return nil, err
		}

		// Calculate point
		p, err := NewPoint(PointParameters{
			Candlesticks: candlesticks,
			PriceType:    params.PriceType,
		})
		if err != nil {
			return nil, err
		}

		// Add calculated point to timeserie
		ts.Set(start, p.Price)
	}

	return ts, nil
}

// InvalidValues returns true if there is at least one invalid value in the timeserie.
func InvalidValues(ts *timeserie.TimeSerie[float64]) bool {
	invalidValuesDetected := false
	_ = ts.Loop(func(_ time.Time, v float64) (bool, error) {
		if v == 0 {
			invalidValuesDetected = true
			return true, nil
		}
		return false, nil
	})
	return invalidValuesDetected
}
