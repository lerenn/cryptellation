package event

import (
	"encoding/json"
	"time"

	"github.com/digital-feather/cryptellation/services/backtests/pkg/client/proto"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/status"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/tick"
)

type Event struct {
	Type    Type
	Time    time.Time
	Content interface{}
}

func (e Event) MarshalBinary() ([]byte, error) {
	return json.Marshal(e)
}

func (e *Event) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, e)
}

func (e *Event) ToProtoBuff() (*proto.BacktestEvent, error) {
	content, err := json.Marshal(e.Content)
	if err != nil {
		return nil, err
	}

	return &proto.BacktestEvent{
		Time:    e.Time.Format(time.RFC3339Nano),
		Type:    e.Type.String(),
		Content: string(content),
	}, nil
}

func FromProtoBuff(pb *proto.BacktestEvent) (Event, error) {
	t, err := time.Parse(time.RFC3339Nano, pb.Time)
	if err != nil {
		return Event{}, err
	}

	switch pb.Type {
	case TypeIsTick.String():
		content, err := tick.FromJSON([]byte(pb.Content))
		if err != nil {
			return Event{}, err
		}

		return Event{
			Time:    t,
			Type:    Type(pb.Type),
			Content: content,
		}, nil
	case TypeIsStatus.String():
		content, err := status.FromJSON([]byte(pb.Content))
		if err != nil {
			return Event{}, err
		}

		return Event{
			Time:    t,
			Type:    Type(pb.Type),
			Content: content,
		}, nil
	default:
		return Event{}, ErrUnknownType
	}
}

func OnlyKeepEarliestSameTimeEvents(evts []Event, endTime time.Time) (earliestTime time.Time, filtered []Event) {
	earliestTime = endTime
	for _, e := range evts {
		if earliestTime.After(e.Time) {
			earliestTime = e.Time
			filtered = make([]Event, 0, len(evts))
		}

		if earliestTime.Equal(e.Time) {
			filtered = append(filtered, e)
		}
	}

	return earliestTime, filtered
}
