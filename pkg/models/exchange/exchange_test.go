package exchange

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestExchangeSuite(t *testing.T) {
	suite.Run(t, new(ExchangeSuite))
}

type ExchangeSuite struct {
	suite.Suite
}

//nolint:funlen
func (suite *ExchangeSuite) TestEquals() {
	cases := []struct {
		Exchange1 Exchange
		Exchange2 Exchange
		Expected  bool
	}{
		// Equals
		{
			Exchange1: Exchange{
				Name:    "exchange1",
				Pairs:   []string{"ABC-DEF", "DEF-XYZ"},
				Periods: []string{"M1", "M15"},
				Fees:    0.1,
			},
			Exchange2: Exchange{
				Name:    "exchange1",
				Pairs:   []string{"DEF-XYZ", "ABC-DEF"},
				Periods: []string{"M15", "M1"},
				Fees:    0.1,
			},
			Expected: true,
		},
		// Not equals
		{
			Exchange1: Exchange{
				Name:    "exchange1",
				Pairs:   []string{"ABC-DEF", "DEF-XYZ"},
				Periods: []string{"M1", "M15"},
				Fees:    0.1,
			},
			Exchange2: Exchange{
				Name:    "exchange2",
				Pairs:   []string{"ABC-DEF"},
				Periods: []string{"M1"},
				Fees:    0.1,
			},
			Expected: false,
		},
		{
			Exchange1: Exchange{
				Name:    "exchange1",
				Pairs:   []string{"ABC-DEF", "DEF-XYZ"},
				Periods: []string{"M1", "M15"},
				Fees:    0.1,
			},
			Exchange2: Exchange{
				Name:    "exchange1",
				Pairs:   []string{"ABC-DEF"},
				Periods: []string{"M15", "M1"},
				Fees:    0.1,
			},
			Expected: false,
		},
	}

	for i, c := range cases {
		if c.Expected != c.Exchange1.Equals(c.Exchange2) {
			suite.Require().Failf(
				"Difference with expectation", "case %d is %+v, should be %+v",
				i, c.Expected, c.Exchange1.Equals(c.Exchange2))
		}
	}
}

func (suite *ExchangeSuite) TestMerge() {
	cases := []struct {
		Exchange1 Exchange
		Exchange2 Exchange
		Expected  Exchange
	}{
		{
			Exchange1: Exchange{
				Name:    "exchange1",
				Pairs:   []string{"ABC-DEF", "DEF-XYZ"},
				Periods: []string{"M1", "M15"},
				Fees:    0.1,
			},
			Exchange2: Exchange{
				Name:    "exchange2",
				Pairs:   []string{"ABC-DEF", "ABC-XYZ"},
				Periods: []string{"M1", "M3"},
				Fees:    0.2,
			},
			Expected: Exchange{
				Name:    "exchange1",
				Pairs:   []string{"ABC-DEF", "DEF-XYZ", "ABC-XYZ"},
				Periods: []string{"M1", "M15", "M3"},
				Fees:    0.1,
			},
		},
		{
			Exchange1: Exchange{
				Name:    "exchange1",
				Pairs:   []string{"ABC-DEF", "DEF-XYZ"},
				Periods: []string{"M1", "M15"},
				Fees:    0.1,
			},
			Exchange2: Exchange{
				Name:    "exchange1",
				Pairs:   []string{"ABC-DEF", "DEF-XYZ"},
				Periods: []string{"M1", "M15"},
				Fees:    0.1,
			},
			Expected: Exchange{
				Name:    "exchange1",
				Pairs:   []string{"ABC-DEF", "DEF-XYZ"},
				Periods: []string{"M1", "M15"},
				Fees:    0.1,
			},
		},
	}

	for i, c := range cases {
		merged := c.Exchange1.Merge(c.Exchange2)
		if !c.Expected.Equals(merged) {
			suite.Require().Failf(
				"Difference with expectation", "case %d is %+v, should be %+v", i, c.Expected, merged)
		}
	}
}
