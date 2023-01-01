package exchange

import (
	"reflect"
	"testing"
	"time"

	"github.com/digital-feather/cryptellation/services/exchanges/clients/go/proto"
	"github.com/stretchr/testify/suite"
)

func TestExchangeSuite(t *testing.T) {
	suite.Run(t, new(ExchangeSuite))
}

type ExchangeSuite struct {
	suite.Suite
}

func (suite *ExchangeSuite) TestMerge() {
	cases := []struct {
		Exchange1 Exchange
		Exchange2 Exchange
		Expected  Exchange
	}{
		{
			Exchange1: Exchange{
				Name:           "exchange1",
				PairsSymbols:   []string{"ABC-DEF", "DEF-XYZ"},
				PeriodsSymbols: []string{"M1", "M15"},
				Fees:           0.1,
			},
			Exchange2: Exchange{
				Name:           "exchange2",
				PairsSymbols:   []string{"ABC-DEF", "ABC-XYZ"},
				PeriodsSymbols: []string{"M1", "M3"},
				Fees:           0.2,
			},
			Expected: Exchange{
				Name:           "exchange1",
				PairsSymbols:   []string{"ABC-DEF", "DEF-XYZ", "ABC-XYZ"},
				PeriodsSymbols: []string{"M1", "M15", "M3"},
				Fees:           0.1,
			},
		},
		{
			Exchange1: Exchange{
				Name:           "exchange1",
				PairsSymbols:   []string{"ABC-DEF", "DEF-XYZ"},
				PeriodsSymbols: []string{"M1", "M15"},
				Fees:           0.1,
			},
			Exchange2: Exchange{
				Name:           "exchange1",
				PairsSymbols:   []string{"ABC-DEF", "DEF-XYZ"},
				PeriodsSymbols: []string{"M1", "M15"},
				Fees:           0.1,
			},
			Expected: Exchange{
				Name:           "exchange1",
				PairsSymbols:   []string{"ABC-DEF", "DEF-XYZ"},
				PeriodsSymbols: []string{"M1", "M15"},
				Fees:           0.1,
			},
		},
	}

	for i, c := range cases {
		merged := c.Exchange1.Merge(c.Exchange2)
		if !reflect.DeepEqual(c.Expected, merged) {
			suite.Require().Fail("Difference with expectation for case %d: %+v", i, merged)
		}
	}
}

func (suite *ExchangeSuite) TestFromProtoBuf() {
	exch, err := FromProtoBuf(&proto.Exchange{
		Name: "exchange",
		Periods: []string{
			"M1",
		},
		Pairs: []string{
			"UTC-USDC",
		},
		Fees:         1,
		LastSyncTime: "1970-01-01T00:00:00Z",
	})
	suite.Require().NoError(err)

	suite.Require().Equal("exchange", exch.Name)
	suite.Require().Equal([]string{"M1"}, exch.PeriodsSymbols)
	suite.Require().Equal([]string{"UTC-USDC"}, exch.PairsSymbols)
	suite.Require().Equal(1.0, exch.Fees)
	suite.Require().WithinDuration(time.Unix(0, 0), exch.LastSyncTime, time.Second)
}

func (suite *ExchangeSuite) TestToProtoBuf() {
	e := Exchange{
		Name:           "exchange",
		PairsSymbols:   []string{"ABC-DEF"},
		PeriodsSymbols: []string{"M1"},
		Fees:           0.1,
		LastSyncTime:   time.Unix(60, 0),
	}

	pb := e.ToProfoBuff()
	suite.Require().Equal("exchange", pb.Name)
	suite.Require().Len(pb.Pairs, 1)
	suite.Require().Equal("ABC-DEF", pb.Pairs[0])
	suite.Require().Equal("M1", pb.Periods[0])
	suite.Require().Len(pb.Periods, 1)
	suite.Require().Equal(float64(0.1), pb.Fees)
	suite.Require().Equal("1970-01-01T00:01:00Z", pb.LastSyncTime)
}
