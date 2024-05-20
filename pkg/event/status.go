package event

import (
	"time"
)

type Status struct {
	Finished bool `json:"finished"`
}

func NewStatusEvent(t time.Time, content Status) Event {
	return Event{
		Type:    TypeIsStatus,
		Time:    t,
		Content: content,
	}
}
