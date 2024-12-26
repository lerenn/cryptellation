package candlestick

import (
	"errors"
	"fmt"
	"time"

	"github.com/lerenn/cryptellation/v1/pkg/models/period"
	"github.com/lerenn/cryptellation/v1/pkg/models/timeserie"
)

var (
	// ErrPeriodMismatch is returned when the period of the candlestick does not
	// match the period of the list.
	ErrPeriodMismatch = errors.New("period-mismatch")
	// ErrCandlestickType is returned when the struct is not a candlestick.
	ErrCandlestickType = errors.New("struct-not-candlestick")
	// ErrExchangeMismatch is returned when the exchange of the list does not
	// match the exchange of the list to merge.
	ErrExchangeMismatch = errors.New("exchange-mismatch")
	// ErrPairMismatch is returned when the pair of the list does not match the
	// pair of the list to merge.
	ErrPairMismatch = errors.New("pair-mismatch")
)

// ListMetadata is the metadata of a list of candlesticks.
type ListMetadata struct {
	Exchange string
	Pair     string
	Period   period.Symbol
}

// List is a list of candlesticks with some metadata.
type List struct {
	Metadata ListMetadata
	Data     timeserie.TimeSerie[Candlestick]
}

// NewList creates a new list of candlesticks with the given metadata.
func NewList(exchange, pair string, period period.Symbol) *List {
	return &List{
		Metadata: ListMetadata{
			Exchange: exchange,
			Pair:     pair,
			Period:   period,
		},
		Data: *timeserie.New[Candlestick](),
	}
}

// NewListWithMetadata creates a new list of candlesticks with the same metadata as the
// given list.
func NewListWithMetadata(md ListMetadata) *List {
	return NewList(md.Exchange, md.Pair, md.Period)
}

// FillMissing will add the 'filling' candlestick at each interval between
// 'start' included and 'end' included when there is a missing candlestick at
// the tested interval.
func (l *List) FillMissing(start, end time.Time, filling Candlestick) error {
	d := l.Metadata.Period.Duration()
	for current := start; current.Before(end.Add(d)); current = current.Add(d) {
		// Check if the candlestick exists
		_, exists := l.Data.Get(current)
		if exists {
			continue
		}

		// Set with filling
		filling.Time = current
		if err := l.Set(filling); err != nil {
			return err
		}
	}

	return nil
}

// MustSet will set the candlestick in the list and panic if there is an error.
func (l *List) MustSet(c Candlestick) *List {
	err := l.Set(c)
	if err != nil {
		panic(err)
	}
	return l
}

// Set will set the candlestick in the list.
func (l *List) Set(c Candlestick) error {
	if !l.Metadata.Period.IsAligned(c.Time) {
		return ErrPeriodMismatch
	}

	l.Data.Set(c.Time, c)
	return nil
}

// Merge will merge the given list into the current list.
func (l *List) Merge(l2 *List, options *timeserie.MergeOptions) error {
	switch {
	case l.Metadata.Exchange != l2.Metadata.Exchange:
		return ErrExchangeMismatch
	case l.Metadata.Pair != l2.Metadata.Pair:
		return ErrPairMismatch
	case l.Metadata.Period != l2.Metadata.Period:
		return ErrPeriodMismatch
	default:
		return l.Data.Merge(l2.Data, options)
	}
}

// Last will return the last candlestick of the list.
func (l List) Last() (Candlestick, bool) {
	_, cs, err := l.Data.Last()
	return cs, err
}

// Extract will extract a new list from the current list with the given time
// range and limit.
func (l List) Extract(start, end time.Time, limit uint) *List {
	el := NewList(l.Metadata.Exchange, l.Metadata.Pair, l.Metadata.Period)
	el.Data = *l.Data.Extract(start, end, int(limit))
	return el
}

// ReplaceUncomplete will replace the uncomplete candlesticks of the current list
// with the candlesticks of the given list.
func (l *List) ReplaceUncomplete(l2 *List) {
	_ = l.Loop(func(cs Candlestick) (bool, error) {
		if cs.Uncomplete {
			ucs, exists := l2.Data.Get(cs.Time)
			if exists {
				l.Data.Set(cs.Time, ucs)
			}
		}
		return false, nil
	})
}

// Loop will loop over the candlesticks of the list and call the given function
// for each candlestick. If the function returns true, the loop will stop.
func (l *List) Loop(f func(Candlestick) (bool, error)) error {
	return l.Data.Loop(func(_ time.Time, cs Candlestick) (bool, error) {
		return f(cs)
	})
}

// MergeListIntoOneCandlestick will merge all the candlesticks of the list into
// one candlestick. The time of the candlestick will be the first time of the
// list, rounded to the period. The high will be the highest high, the low the
// lowest low, the volume the sum of all volumes and the close the close of the
// last candlestick.
func MergeListIntoOneCandlestick(csl *List, per period.Symbol) (time.Time, Candlestick) {
	if csl.Data.Len() == 0 {
		return time.Unix(0, 0), Candlestick{}
	}

	mts, mcs, _ := csl.Data.First()
	mts = per.RoundTime(mts)

	_ = csl.Loop(func(cs Candlestick) (bool, error) {
		if !per.RoundTime(cs.Time).Equal(mts) {
			return true, nil
		}

		if cs.Time.Equal(mts) {
			return false, nil
		}

		if cs.High > mcs.High {
			mcs.High = cs.High
		}
		if cs.Low < mcs.Low {
			mcs.Low = cs.Low
		}
		mcs.Volume += cs.Volume
		mcs.Close = cs.Close

		return false, nil
	})

	return mts, mcs
}

// String returns a string representation of the list.
func (l List) String() string {
	txt := fmt.Sprintf("# %s - %s - %s\n", l.Metadata.Exchange, l.Metadata.Pair, l.Metadata.Period.String())

	_ = l.Loop(func(cs Candlestick) (bool, error) {
		uncomplete := ""
		if cs.Uncomplete {
			uncomplete = "uncomplete"
		}

		txt += fmt.Sprintf(
			" %s: %04f/%04f/%04f/%04f (%f) %s\n",
			cs.Time.Format(time.RFC3339),
			cs.Open,
			cs.High,
			cs.Low,
			cs.Close,
			cs.Volume,
			uncomplete)
		return false, nil
	})

	return txt
}

// GetUncompleteTimes returns an array of time from candlesticks that are marked
// as uncomplete (i.e. data pulled when candlestick covering time was not complete).
func (l List) GetUncompleteTimes() []time.Time {
	uncomplete := make([]time.Time, 0)
	_ = l.Loop(func(cs Candlestick) (bool, error) {
		if cs.Uncomplete {
			uncomplete = append(uncomplete, cs.Time)
		}
		return false, nil
	})
	return uncomplete
}

// GetUncompleteRange returns an array of time range from candlesticks that are
// marked as uncomplete (i.e. data pulled when candlestick covering time was not
// complete).
func (l List) GetUncompleteRange() []timeserie.TimeRange {
	// Change to ranges
	ut := l.GetUncompleteTimes()
	tr := make([]timeserie.TimeRange, len(ut))
	for i, t := range ut {
		tr[i].Start, tr[i].End = t, t
	}

	// Merge everything
	tr, _ = timeserie.MergeTimeRanges(tr, tr)
	return tr
}

// GetMissingTimes returns an array of missing time in the candlestick list.
func (l List) GetMissingTimes(start, end time.Time, limit uint) []time.Time {
	// Quick check that there is missing
	if !l.Data.AreMissing(start, end, l.Metadata.Period.Duration(), limit) {
		return []time.Time{}
	}

	// Get missing times from timeserie
	return l.Data.GetMissingTimes(start, end, l.Metadata.Period.Duration(), limit)
}

// GetMissingRange returns an array of missing time range in the candlestick list.
func (l List) GetMissingRange(start, end time.Time, limit uint) []timeserie.TimeRange {
	// Quick check that there is missing
	if !l.Data.AreMissing(start, end, l.Metadata.Period.Duration(), limit) {
		return []timeserie.TimeRange{}
	}

	// Get missing range from timeserie
	return l.Data.GetMissingRanges(start, end, l.Metadata.Period.Duration(), limit)
}

// ToArray returns an array of candlesticks from the list.
func (l List) ToArray() []Candlestick {
	cs := make([]Candlestick, 0, l.Data.Len())
	_ = l.Loop(func(c Candlestick) (bool, error) {
		cs = append(cs, c)
		return false, nil
	})
	return cs
}
