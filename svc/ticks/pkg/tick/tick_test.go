package tick

import (
	"encoding/json"
	"testing"
	"time"

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
		Time:     time.Unix(60, 0).UTC(),
		Pair:     "BTC-USDC",
		Price:    1.01,
		Exchange: "exchange",
	}

	b, err := json.Marshal(tick)
	suite.Require().NoError(err)

	tick2 := Tick{}
	suite.Require().NoError(json.Unmarshal(b, &tick2))
	suite.Require().Equal(tick, tick2)
}
