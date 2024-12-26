package sma

import (
	"fmt"
	"time"

	"github.com/lerenn/cryptellation/v1/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/v1/pkg/models/period"
)

// PointParameters is the params for the Point function.
type PointParameters struct {
	Candlesticks *candlestick.List
	PriceType    candlestick.PriceType
}

// Point is a point of the SMA.
type Point struct {
	Exchange  string
	Pair      string
	Period    period.Symbol
	PeriodNb  int
	PriceType candlestick.PriceType
	Time      time.Time
	Price     float64
}

// NewPoint creates a new point from the given parameters.
func NewPoint(params PointParameters) (Point, error) {
	var total float64

	// Get count of candlesticks
	count := params.Candlesticks.Data.Len()

	// Get total from the timeserie
	if err := params.Candlesticks.Loop(func(cs candlestick.Candlestick) (bool, error) {
		price := cs.Price(params.PriceType)

		// Reduce the count if the price is 0
		if price == 0 {
			count--
			return false, nil
		}

		total += price

		return false, nil
	}); err != nil {
		return Point{}, err
	}

	// Get point time
	last, ok := params.Candlesticks.Last()
	if !ok {
		return Point{}, fmt.Errorf("no last candlestick")
	}

	// Get average price
	price := total / float64(count)

	return Point{
		Exchange:  params.Candlesticks.Metadata.Exchange,
		Pair:      params.Candlesticks.Metadata.Pair,
		Period:    params.Candlesticks.Metadata.Period,
		PeriodNb:  params.Candlesticks.Data.Len(),
		PriceType: params.PriceType,
		Time:      last.Time,
		Price:     price,
	}, nil
}
