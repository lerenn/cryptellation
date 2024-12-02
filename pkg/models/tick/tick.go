package tick

import (
	"encoding/json"
	"fmt"
	"time"
)

type Tick struct {
	Time     time.Time `json:"time"`
	Pair     string    `json:"pair"`
	Price    float64   `json:"price"`
	Exchange string    `json:"exchange"`
}

func (t Tick) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Tick) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, t)
}

func (t Tick) String() string {
	return fmt.Sprintf(
		"[%s|%s|%s] %f",
		t.Time.Format(time.RFC3339),
		t.Exchange,
		t.Pair,
		t.Price,
	)
}
