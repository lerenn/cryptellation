package event

import (
	"testing"
	"time"

	"github.com/digital-feather/cryptellation/services/backtests/pkg/client/proto"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/status"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/tick"
	"github.com/stretchr/testify/suite"
)

func TestEventSuite(t *testing.T) {
	suite.Run(t, new(EventSuite))
}

type EventSuite struct {
	suite.Suite
}

func (suite *EventSuite) TestOnlyKeepEarliestSameTimeEvents() {
	cases := []struct {
		In      []Event
		InTime  time.Time
		Out     []Event
		OutTime time.Time
	}{
		{
			In:      []Event{},
			InTime:  time.Unix(1<<62, 0),
			Out:     []Event{},
			OutTime: time.Unix(1<<62, 0),
		},
		{
			In: []Event{
				NewTickEvent(time.Unix(120, 0), tick.Tick{}),
				NewTickEvent(time.Unix(60, 0), tick.Tick{}),
				NewTickEvent(time.Unix(240, 0), tick.Tick{}),
				NewTickEvent(time.Unix(60, 0), tick.Tick{}),
				NewTickEvent(time.Unix(180, 0), tick.Tick{}),
			},
			InTime: time.Unix(1<<62, 0),
			Out: []Event{
				NewTickEvent(time.Unix(60, 0), tick.Tick{}),
				NewTickEvent(time.Unix(60, 0), tick.Tick{}),
			},
			OutTime: time.Unix(60, 0),
		},
	}

	for i, c := range cases {
		t, out := OnlyKeepEarliestSameTimeEvents(c.In, c.InTime)
		suite.Require().WithinDuration(c.OutTime, t, time.Microsecond, i)
		suite.Require().Len(out, len(c.Out), i)
	}
}

func (suite *EventSuite) TestToProtoBuffWithTickEvent() {
	evt := Event{
		Type: TypeIsTick,
		Time: time.Unix(60, 0),
		Content: tick.Tick{
			PairSymbol: "BTC-USDC",
			Price:      1.01,
			Exchange:   "exchange",
		},
	}

	pb, err := evt.ToProtoBuff()
	suite.NoError(err)
	suite.Require().Equal(evt.Time.Format(time.RFC3339Nano), pb.Time)
	suite.Require().Equal(evt.Type.String(), pb.Type)
	suite.Require().Equal("{\"pair_symbol\":\"BTC-USDC\",\"price\":1.01,\"exchange\":\"exchange\"}", pb.Content)
}

func (suite *EventSuite) TestFromProtoBuffWithTickEvent() {
	pbTick := &proto.BacktestEvent{
		Time:    "1970-01-01T00:01:00Z",
		Type:    TypeIsTick.String(),
		Content: "{\"pair_symbol\":\"BTC-USDC\",\"price\":1.01,\"exchange\":\"exchange\"}",
	}

	expectedTick := tick.Tick{
		PairSymbol: "BTC-USDC",
		Price:      1.01,
		Exchange:   "exchange",
	}

	t, err := FromProtoBuff(pbTick)
	suite.Require().NoError(err)
	suite.Require().WithinDuration(time.Unix(60, 0).UTC(), t.Time, time.Millisecond)
	suite.Require().Equal(TypeIsTick, t.Type)
	suite.Require().Equal(expectedTick, t.Content)
}

func (suite *EventSuite) TestToProtoBuffWithStatusEvent() {
	evt := Event{
		Type: TypeIsTick,
		Time: time.Unix(60, 0),
		Content: status.Status{
			Finished: true,
		},
	}

	pb, err := evt.ToProtoBuff()
	suite.NoError(err)
	suite.Require().Equal(evt.Time.Format(time.RFC3339Nano), pb.Time)
	suite.Require().Equal(evt.Type.String(), pb.Type)
	suite.Require().Equal("{\"finished\":true}", pb.Content)
}

func (suite *EventSuite) TestFromProtoBuffWithStatusEvent() {
	pbTick := &proto.BacktestEvent{
		Time:    "1970-01-01T00:01:00Z",
		Type:    TypeIsStatus.String(),
		Content: "{\"finished\":true}",
	}

	expectedTick := status.Status{
		Finished: true,
	}

	t, err := FromProtoBuff(pbTick)
	suite.Require().NoError(err)
	suite.Require().WithinDuration(time.Unix(60, 0).UTC(), t.Time, time.Millisecond)
	suite.Require().Equal(TypeIsStatus, t.Type)
	suite.Require().Equal(expectedTick, t.Content)
}
