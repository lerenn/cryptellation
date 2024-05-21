package sma

import (
	"time"

	"github.com/lerenn/cryptellation/pkg/models/timeserie"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
)

type PointPayload struct {
	Candlesticks *timeserie.TimeSerie[candlestick.Candlestick]
	PriceType    candlestick.PriceType
}

func Point(payload PointPayload) float64 {
	var total float64

	// Get total from the timeserie
	_ = payload.Candlesticks.Loop(func(t time.Time, cs candlestick.Candlestick) (bool, error) {
		total += cs.PriceByType(payload.PriceType)
		return false, nil
	})

	// Get average
	return total / float64(payload.Candlesticks.Len())
}
