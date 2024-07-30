package candlesticks

import (
	"math"
	"time"

	"cryptellation/svc/candlesticks/pkg/candlestick"
)

func getMinMax(data *candlestick.List) (min, max float64) {
	min, max = math.MaxFloat64, 0
	_ = data.Loop(func(t time.Time, c candlestick.Candlestick) (bool, error) {
		if c.Low < min {
			min = c.Low
		}
		if c.High > max {
			max = c.High
		}

		return false, nil
	})
	return
}
