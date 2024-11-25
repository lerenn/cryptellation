package candlestick

import (
	"fmt"
	"time"
)

// Candlestick is a candlestick with the time, open, high, low, close and volume.
// Uncomplete is true if the candlestick is not closed yet or have not been
// updated since it was closed.
type Candlestick struct {
	Time       time.Time `bson:"time"    json:"time,omitempty"`
	Open       float64   `bson:"open"     json:"open,omitempty"`
	High       float64   `bson:"high"     json:"high,omitempty"`
	Low        float64   `bson:"low"      json:"low,omitempty"`
	Close      float64   `bson:"close"    json:"close,omitempty"`
	Volume     float64   `bson:"volume"   json:"volume,omitempty"`
	Uncomplete bool      `bson:"complete" json:"uncomplete,omitempty"`
}

// Equal checks equality between the candlesticks.
func (cs Candlestick) Equal(b Candlestick) bool {
	t := cs.Time.Equal(b.Time)
	o := cs.Open == b.Open
	h := cs.High == b.High
	l := cs.Low == b.Low
	c := cs.Close == b.Close
	v := cs.Volume == b.Volume
	u := cs.Uncomplete == b.Uncomplete
	return t && o && h && l && c && v && u
}

// Price is the price of the candlestick depending on the price type.
func (cs Candlestick) Price(p Price) float64 {
	switch p {
	case PriceIsOpen:
		return cs.Open
	case PriceIsHigh:
		return cs.High
	case PriceIsLow:
		return cs.Low
	case PriceIsClose:
		fallthrough
	default:
		return cs.Close
	}
}

// String is a string representation of the candlestick.
func (cs Candlestick) String() string {
	return fmt.Sprintf("T: %s | O:%f | H:%f | L:%f | C:%f", cs.Time, cs.Open, cs.High, cs.Low, cs.Close)
}
