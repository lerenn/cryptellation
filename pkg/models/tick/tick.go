package tick

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Tick is the struct that will handle the ticks.
type Tick struct {
	Time     time.Time `json:"time"`
	Pair     string    `json:"pair"`
	Price    float64   `json:"price"`
	Exchange string    `json:"exchange"`
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
