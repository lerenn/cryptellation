package timeserie

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

func TestTimeRangeSuite(t *testing.T) {
	suite.Run(t, new(TimeRangeSuite))
}

type TimeRangeSuite struct {
	suite.Suite
}

func (suite *TimeRangeSuite) TestMergeTimeRanges() {
	cases := []struct {
		tr1, tr2, expected []TimeRange
	}{
		// Full empty
		{
			tr1:      []TimeRange{},
			tr2:      []TimeRange{},
			expected: []TimeRange{},
		},
		// TR1 empty
		{
			tr1:      []TimeRange{},
			tr2:      []TimeRange{{Start: time.Unix(60, 0), End: time.Unix(240, 0)}},
			expected: []TimeRange{{Start: time.Unix(60, 0), End: time.Unix(240, 0)}},
		},
		// TR2 empty
		{
			tr1:      []TimeRange{{Start: time.Unix(0, 0), End: time.Unix(180, 0)}},
			tr2:      []TimeRange{},
			expected: []TimeRange{{Start: time.Unix(0, 0), End: time.Unix(180, 0)}},
		},
		// Overlapping
		{
			tr1:      []TimeRange{{Start: time.Unix(0, 0), End: time.Unix(180, 0)}},
			tr2:      []TimeRange{{Start: time.Unix(60, 0), End: time.Unix(240, 0)}},
			expected: []TimeRange{{Start: time.Unix(0, 0), End: time.Unix(240, 0)}},
		},
		// Following
		{
			tr1:      []TimeRange{{Start: time.Unix(0, 0), End: time.Unix(180, 0)}},
			tr2:      []TimeRange{{Start: time.Unix(180, 0), End: time.Unix(240, 0)}},
			expected: []TimeRange{{Start: time.Unix(0, 0), End: time.Unix(240, 0)}},
		},
		// TR1 == TR2
		{
			tr1:      []TimeRange{{Start: time.Unix(0, 0), End: time.Unix(180, 0)}},
			tr2:      []TimeRange{{Start: time.Unix(0, 0), End: time.Unix(180, 0)}},
			expected: []TimeRange{{Start: time.Unix(0, 0), End: time.Unix(180, 0)}},
		},
		// Unordered TR1
		{
			tr1: []TimeRange{
				{Start: time.Unix(300, 0), End: time.Unix(360, 0)},
				{Start: time.Unix(0, 0), End: time.Unix(180, 0)},
			},
			tr2: []TimeRange{{Start: time.Unix(180, 0), End: time.Unix(240, 0)}},
			expected: []TimeRange{
				{Start: time.Unix(0, 0), End: time.Unix(240, 0)},
				{Start: time.Unix(300, 0), End: time.Unix(360, 0)},
			},
		},
	}

	for i, c := range cases {
		res, err := MergeTimeRanges(c.tr1, c.tr2)
		suite.Require().NoError(err)
		suite.Require().Equal(c.expected, res, i)
	}
}

func (suite *TimeRangeSuite) TestTimeRangesFromMissingTimes() {
	cases := []struct {
		interval time.Duration
		times    []time.Time
		expected []TimeRange
	}{
		// Empty
		{
			interval: time.Minute,
			times:    []time.Time{},
			expected: []TimeRange{},
		},
		// Single
		{
			interval: time.Minute,
			times:    []time.Time{time.Unix(60, 0)},
			expected: []TimeRange{{Start: time.Unix(60, 0), End: time.Unix(60, 0)}},
		},
		// Multiple only consecutive
		{
			interval: time.Minute,
			times: []time.Time{
				time.Unix(60, 0),
				time.Unix(120, 0),
				time.Unix(180, 0),
			},
			expected: []TimeRange{
				{Start: time.Unix(60, 0), End: time.Unix(180, 0)},
			},
		},
		// Multiple with non-consecutive
		{
			interval: time.Minute,
			times: []time.Time{
				time.Unix(60, 0),
				time.Unix(120, 0),
				time.Unix(240, 0),
			},
			expected: []TimeRange{
				{Start: time.Unix(60, 0), End: time.Unix(120, 0)},
				{Start: time.Unix(240, 0), End: time.Unix(240, 0)},
			},
		},
	}

	for i, c := range cases {
		res := TimeRangesFromMissingTimes(c.interval, c.times)
		suite.Require().Equal(c.expected, res, i)
	}
}
