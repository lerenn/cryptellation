package tick

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/digital-feather/cryptellation/services/ticks/pkg/client/proto"
	"github.com/stretchr/testify/suite"
)

func TestTickSuite(t *testing.T) {
	suite.Run(t, new(TickSuite))
}

type TickSuite struct {
	suite.Suite
}

func (suite *TickSuite) TestMarshalingJSON() {
	tick := Tick{
		Time:       time.Unix(60, 0).UTC(),
		PairSymbol: "BTC-USDC",
		Price:      1.01,
		Exchange:   "exchange",
	}

	b, err := json.Marshal(tick)
	suite.Require().NoError(err)

	tick2 := Tick{}
	suite.Require().NoError(json.Unmarshal(b, &tick2))
	suite.Require().Equal(tick, tick2)
}

func (suite *TickSuite) TestToProtoBuff() {
	tick := Tick{
		Time:       time.Unix(60, 0).UTC(),
		PairSymbol: "BTC-USDC",
		Price:      1.01,
		Exchange:   "exchange",
	}

	pb := tick.ToProtoBuff()
	suite.Require().Equal(tick.Time.Format(time.RFC3339Nano), pb.Time)
	suite.Require().Equal(tick.PairSymbol, pb.PairSymbol)
	suite.Require().Equal(float32(tick.Price), pb.Price)
	suite.Require().Equal(tick.Exchange, pb.Exchange)
}

func (suite *TickSuite) TestFromProtoBuff() {
	pbTick := &proto.Tick{
		Time:       "1970-01-01T00:01:00Z",
		PairSymbol: "BTC-USDC",
		Price:      1.01,
		Exchange:   "exchange",
	}

	t, err := FromProtoBuff(pbTick)
	suite.Require().NoError(err)
	suite.Require().WithinDuration(time.Unix(60, 0).UTC(), t.Time, time.Millisecond)
	suite.Require().Equal(pbTick.PairSymbol, t.PairSymbol)
	suite.Require().Equal(pbTick.Price, float32(t.Price))
	suite.Require().Equal(pbTick.Exchange, t.Exchange)
}
