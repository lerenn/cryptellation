package tick

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/lerenn/cryptellation/v1/pkg/models/candlestick"
)

// Tick is the struct that will handle the ticks.
type Tick struct {
	Time     time.Time `json:"time"`
	Pair     string    `json:"pair"`
	Price    float64   `json:"price"`
	Exchange string    `json:"exchange"`
}

// FromCandlestick creates a Tick from a candlestick.
func FromCandlestick(
	exchange, pair string,
	currentPriceType candlestick.Price,
	t time.Time,
	cs candlestick.Candlestick,
) Tick {
	return Tick{
		Time:     t,
		Pair:     pair,
		Price:    cs.Price(currentPriceType),
		Exchange: exchange,
	}
}

// MarshalBinary marshals a Tick into a byte slice.
func (t Tick) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

// UnmarshalBinary unmarshals a byte slice into a Tick.
func (t *Tick) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, t)
}

// String returns a string representation of a Tick.
func (t Tick) String() string {
	return fmt.Sprintf(
		"[ %-30s | %s | %s ] %f",
		t.Time.Format(time.RFC3339Nano),
		strings.ToTitle(t.Exchange),
		t.Pair,
		t.Price,
	)
}

// OnlyKeepEarliestSameTime returns the earliest time and the ticks that have the same time as the earliest time.
func OnlyKeepEarliestSameTime(originals []Tick, endTime time.Time) (earliestTime time.Time, filtered []Tick) {
	earliestTime = endTime
	for _, e := range originals {
		if earliestTime.After(e.Time) {
			earliestTime = e.Time
			filtered = make([]Tick, 0, len(originals))
		}

		if earliestTime.Equal(e.Time) {
			filtered = append(filtered, e)
		}
	}

	return earliestTime, filtered
}
