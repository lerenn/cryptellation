package sma

import (
	"time"

	"github.com/lerenn/cryptellation/pkg/models/timeserie"

	"github.com/lerenn/cryptellation/candlesticks/pkg/candlestick"
)

type PointPayload struct {
	Candlesticks *timeserie.TimeSerie[candlestick.Candlestick]
	PriceType    candlestick.PriceType
}

func Point(payload PointPayload) float64 {
	var total float64

	// Get count of candlesticks
	count := payload.Candlesticks.Len()

	// Get total from the timeserie
	_ = payload.Candlesticks.Loop(func(t time.Time, cs candlestick.Candlestick) (bool, error) {
		price := cs.PriceByType(payload.PriceType)

		// Reduce the count if the price is 0
		if price == 0 {
			count--
			return false, nil
		}

		total += price

		return false, nil
	})

	// Get average
	return total / float64(count)
}
