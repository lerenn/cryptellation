package tick

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestTickSuite(t *testing.T) {
	suite.Run(t, new(TickSuite))
}

type TickSuite struct {
	suite.Suite
}

func (suite *TickSuite) TestFromJSON() {
	tests := []struct {
		Input       []byte
		Output      Tick
		OutputError bool
	}{
		{
			Input:       []byte("{\"pair_symbol\":\"ETH-USDC\",\"price\":1.2,\"exchange\":\"exchange\"}"),
			Output:      Tick{PairSymbol: "ETH-USDC", Price: 1.2, Exchange: "exchange"},
			OutputError: false,
		},
		{
			Input:       []byte("{\"pair_symbol\":\"ETH-USDC\",\"price\":\"1.2\",\"exchange\":\"exchange\"}"),
			Output:      Tick{},
			OutputError: true,
		},
	}

	for i, t := range tests {
		output, err := FromJSON(t.Input)
		if t.OutputError {
			suite.Require().Error(err, i)
			continue
		} else {
			suite.Require().NoError(err, i)
		}

		suite.Require().Equal(t.Output, output, i)
	}
}

func (suite *TickSuite) TestMarshalingJSON() {
	as := suite.Require()

	tick := Tick{
		PairSymbol: "BTC-USDC",
		Price:      1.01,
		Exchange:   "exchange",
	}

	b, err := json.Marshal(tick)
	as.NoError(err)

	tick2 := Tick{}
	as.NoError(json.Unmarshal(b, &tick2))
	as.Equal(tick, tick2)
}
