package candlesticks

import (
	"math"
	"time"
)

type Candlestick struct {
	Time  time.Time
	Open  float64
	High  float64
	Low   float64
	Close float64
}

func getMinMax(data []Candlestick) (min float64, max float64) {
	min, max = math.MaxFloat64, 0
	for _, d := range data {
		if d.Low < min {
			min = d.Low
		}
		if d.High > max {
			max = d.High
		}
	}
	return
}
