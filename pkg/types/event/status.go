package event

import (
	"encoding/json"
	"time"
)

type Status struct {
	Finished bool `json:"finished"`
}

func FromJSON(content []byte) (Status, error) {
	var st Status
	err := json.Unmarshal(content, &st)
	return st, err
}

func NewStatusEvent(t time.Time, content Status) Event {
	return Event{
		Type:    TypeIsStatus,
		Time:    t,
		Content: content,
	}
}
