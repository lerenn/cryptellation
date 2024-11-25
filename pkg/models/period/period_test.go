package period

import (
	"testing"
	"time"

	"github.com/lerenn/cryptellation/v1/pkg/utils"
	"github.com/stretchr/testify/suite"
)

func TestPeriodSuite(t *testing.T) {
	suite.Run(t, new(PeriodSuite))
}

type PeriodSuite struct {
	suite.Suite
}

func (suite *PeriodSuite) TestPeriodDuration() {
	suite.Require().Equal(time.Minute, M1.Duration())
}

func (suite *PeriodSuite) TestRoundTimeNano() {
	toCorrect := utils.Must(time.Parse(time.RFC3339, "2024-05-29T14:52:00.660748646Z"))
	expected := utils.Must(time.Parse(time.RFC3339, "2024-05-29T14:52:00Z"))

	corrected := M1.RoundTime(toCorrect)
	suite.WithinDuration(expected, corrected, time.Nanosecond)
}

func (suite *PeriodSuite) TestPeriods() {
	symbols := Symbols()

	suite.Require().Len(symbols, 14)
	suite.Require().Contains(symbols, M1)
	suite.Require().Contains(symbols, M3)
	suite.Require().Contains(symbols, M5)
	suite.Require().Contains(symbols, M15)
	suite.Require().Contains(symbols, M30)
	suite.Require().Contains(symbols, H1)
	suite.Require().Contains(symbols, H2)
	suite.Require().Contains(symbols, H4)
	suite.Require().Contains(symbols, H6)
	suite.Require().Contains(symbols, H8)
	suite.Require().Contains(symbols, H12)
	suite.Require().Contains(symbols, D1)
	suite.Require().Contains(symbols, D3)
	suite.Require().Contains(symbols, W1)
}

func (suite *PeriodSuite) TestSymbolsString() {
	for _, s := range Symbols() {
		suite.Require().NotEqual(ErrInvalidPeriod.Error(), s.String())
	}
}

func (suite *PeriodSuite) TestValidateSymbol() {
	for _, s := range Symbols() {
		suite.Require().NoError(s.Validate())
	}

	suite.Require().ErrorIs(Symbol("unknown").Validate(), ErrInvalidPeriod)
}

func (suite *PeriodSuite) TestIsAligned() {
	suite.Require().True(M1.IsAligned(time.Unix(60, 0)))
	suite.Require().False(M1.IsAligned(time.Unix(45, 0)))
}

func (suite *PeriodSuite) TestFromSeconds() {
	s, err := FromSeconds(60)
	suite.Require().NoError(err)
	suite.Require().Equal(M1, s)

	_, err = FromSeconds(59)
	suite.Require().Error(err)
}

func (suite *PeriodSuite) TestCountBetweenTimes() {
	symbs := []Symbol{M1, D1}

	now := time.Now()
	for _, s := range symbs {
		suite.Require().Equal(int64(0), s.CountBetweenTimes(now, now))
		suite.Require().Equal(int64(1), s.CountBetweenTimes(now.Add(-s.Duration()), now))
		suite.Require().Equal(int64(10), s.CountBetweenTimes(now.Add(-s.Duration()*10), now))
	}
}

func (suite *PeriodSuite) TestUniqueArray() {
	s1 := []Symbol{M1, M15}
	s2 := []Symbol{M1, M3}

	m := UniqueArray(s2, s1)
	suite.Require().Len(m, 3)
	suite.Require().Contains(m, M1)
	suite.Require().Contains(m, M3)
	suite.Require().Contains(m, M15)
}

func (suite *PeriodSuite) TestProcessRequestedStartEndTimeWithEmpty() {
	// Given no time

	// When we request the processing
	start, end := M1.RoundInterval(nil, nil)

	// Then we have correct times
	suite.Require().WithinDuration(time.Now().Add(-M1.Duration()*500), start, time.Minute)
	suite.Require().WithinDuration(time.Now(), end, time.Minute)
}

func (suite *PeriodSuite) TestProcessRequestedStartEndTimeWithNoStart() {
	// Given only end
	reqEnd, err := time.Parse(time.RFC3339, "1993-11-15T11:29:01Z")
	suite.Require().NoError(err)

	// When we request the processing
	start, end := M1.RoundInterval(nil, &reqEnd)

	// Then we have correct times
	expectedEnd, err := time.Parse(time.RFC3339, "1993-11-15T11:29:00Z")
	suite.Require().NoError(err)
	expectedStart := expectedEnd.Add(-M1.Duration() * 500)
	suite.Require().WithinDuration(expectedStart, start, time.Millisecond)
	suite.Require().WithinDuration(expectedEnd, end, time.Millisecond)
}

func (suite *PeriodSuite) TestProcessRequestedStartEndTimeWithNoEnd() {
	// Given only start
	reqStart, err := time.Parse(time.RFC3339, "1993-11-15T11:00:01Z")
	suite.Require().NoError(err)

	// When we request the processing
	start, end := M1.RoundInterval(&reqStart, nil)

	// Then we have correct times
	expectedStart, err := time.Parse(time.RFC3339, "1993-11-15T11:00:00Z")
	suite.Require().NoError(err)
	expectedEnd := expectedStart.Add(M1.Duration() * 500)
	suite.Require().WithinDuration(expectedStart, start, time.Millisecond)
	suite.Require().WithinDuration(expectedEnd, end, time.Millisecond)
}

func (suite *PeriodSuite) TestProcessRequestedStartEndTimeWithNonAligned() {
	// Given a time, but not aligned
	reqStart, err := time.Parse(time.RFC3339, "1993-11-15T11:00:04Z")
	suite.Require().NoError(err)
	reqEnd, err := time.Parse(time.RFC3339, "1993-11-15T11:29:01Z")
	suite.Require().NoError(err)

	// When we request the processing
	start, end := M1.RoundInterval(&reqStart, &reqEnd)

	// Then we have correct times
	expectedStart, err := time.Parse(time.RFC3339, "1993-11-15T11:00:00Z")
	suite.Require().NoError(err)
	expectedEnd, err := time.Parse(time.RFC3339, "1993-11-15T11:29:00Z")
	suite.Require().NoError(err)
	suite.Require().WithinDuration(expectedStart, start, time.Millisecond)
	suite.Require().WithinDuration(expectedEnd, end, time.Millisecond)
}
