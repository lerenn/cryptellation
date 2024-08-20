package candlestick

import (
	"errors"
	"fmt"
	"time"

	"github.com/lerenn/cryptellation/pkg/models/timeserie"

	"github.com/lerenn/cryptellation/candlesticks/pkg/period"
)

var (
	ErrPeriodMismatch   = errors.New("period-mismatch")
	ErrCandlestickType  = errors.New("struct-not-candlestick")
	ErrExchangeMismatch = errors.New("exchange-mismatch")
	ErrPairMismatch     = errors.New("pair-mismatch")
)

type List struct {
	Exchange string
	Pair     string
	Period   period.Symbol
	timeserie.TimeSerie[Candlestick]
}

func NewList(exchange, pair string, period period.Symbol) *List {
	return &List{
		Exchange:  exchange,
		Pair:      pair,
		Period:    period,
		TimeSerie: *timeserie.New[Candlestick](),
	}
}

func NewListFrom(l *List) *List {
	return NewList(l.Exchange, l.Pair, l.Period)
}

// FillMissing will add the 'filling' candlestick at each interval between
// 'start' included and 'end' included when there is a missing candlestick at
// the tested interval.
func (l *List) FillMissing(start, end time.Time, filling Candlestick) error {
	for current := start; current.Before(end.Add(l.Period.Duration())); current = current.Add(l.Period.Duration()) {
		// Check if the candlestick exists
		_, exists := l.Get(current)
		if exists {
			continue
		}

		// Set with filling
		if err := l.Set(current, filling); err != nil {
			return err
		}
	}

	return nil
}

func (l *List) MustSet(t time.Time, c Candlestick) *List {
	err := l.Set(t, c)
	if err != nil {
		panic(err)
	}
	return l
}

func (l *List) Set(t time.Time, c Candlestick) error {
	if !l.Period.IsAligned(t) {
		return ErrPeriodMismatch
	}

	l.TimeSerie.Set(t, c)
	return nil
}

func (l *List) Merge(l2 *List, options *timeserie.MergeOptions) error {
	if l.Exchange != l2.Exchange {
		return ErrExchangeMismatch
	} else if l.Pair != l2.Pair {
		return ErrPairMismatch
	} else if l.Period != l2.Period {
		return ErrPeriodMismatch
	}

	return l.TimeSerie.Merge(l2.TimeSerie, options)
}

func (l List) Extract(start, end time.Time, limit uint) *List {
	el := NewList(l.Exchange, l.Pair, l.Period)
	el.TimeSerie = *l.TimeSerie.Extract(start, end, int(limit))
	return el
}

func (l *List) ReplaceUncomplete(l2 *List) {
	_ = l.Loop(func(t time.Time, cs Candlestick) (bool, error) {
		if cs.Uncomplete {
			ucs, exists := l2.Get(t)
			if exists {
				l.TimeSerie.Set(t, ucs)
			}
		}
		return false, nil
	})
}

func MergeListIntoOneCandlestick(csl *List, per period.Symbol) (time.Time, Candlestick) {
	if csl.Len() == 0 {
		return time.Unix(0, 0), Candlestick{}
	}

	mts, mcs, _ := csl.TimeSerie.First()
	mts = per.RoundTime(mts)

	_ = csl.Loop(func(t time.Time, cs Candlestick) (bool, error) {
		if !per.RoundTime(t).Equal(mts) {
			return true, nil
		}

		if t.Equal(mts) {
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

func (l List) String() string {
	txt := fmt.Sprintf("# %s - %s - %s\n", l.Exchange, l.Pair, l.Period.String())

	_ = l.Loop(func(t time.Time, cs Candlestick) (bool, error) {
		uncomplete := ""
		if cs.Uncomplete {
			uncomplete = "uncomplete"
		}

		txt += fmt.Sprintf(
			" %s: %04f/%04f/%04f/%04f (%f) %s\n",
			t.Format(time.RFC3339),
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
// as uncomplete (i.e. data pulled when candlestick covering time was not complete)
func (cl List) GetUncompleteTimes() []time.Time {
	uncomplete := make([]time.Time, 0)
	_ = cl.Loop(func(t time.Time, cs Candlestick) (bool, error) {
		if cs.Uncomplete {
			uncomplete = append(uncomplete, t)
		}
		return false, nil
	})
	return uncomplete
}

func (cl List) GetUncompleteRange() []timeserie.TimeRange {
	// Change to ranges
	ut := cl.GetUncompleteTimes()
	tr := make([]timeserie.TimeRange, len(ut))
	for i, t := range ut {
		tr[i].Start, tr[i].End = t, t
	}

	// Merge everything
	tr, _ = timeserie.MergeTimeRanges(tr, tr)
	return tr
}

// GetMissingTimes returns an array of missing time in the candlestick list
func (cl List) GetMissingTimes(start, end time.Time, limit uint) []time.Time {
	// Quick check that there is missing
	if !cl.TimeSerie.AreMissing(start, end, cl.Period.Duration(), limit) {
		return []time.Time{}
	}

	// Get missing times from timeserie
	return cl.TimeSerie.GetMissingTimes(start, end, cl.Period.Duration(), limit)
}

func (cl List) GetMissingRange(start, end time.Time, limit uint) []timeserie.TimeRange {
	// Quick check that there is missing
	if !cl.TimeSerie.AreMissing(start, end, cl.Period.Duration(), limit) {
		return []timeserie.TimeRange{}
	}

	// Get missing range from timeserie
	return cl.TimeSerie.GetMissingRanges(start, end, cl.Period.Duration(), limit)
}
