package candlestick

import (
	"errors"
	"fmt"
	"time"

	"github.com/lerenn/cryptellation/pkg/models/period"
	"github.com/lerenn/cryptellation/pkg/models/timeserie"
)

var (
	ErrPeriodMismatch   = errors.New("period-mismatch")
	ErrCandlestickType  = errors.New("struct-not-candlestick")
	ErrExchangeMismatch = errors.New("exchange-mismatch")
	ErrPairMismatch     = errors.New("pair-mismatch")
)

type List struct {
	ExchangeName string
	PairSymbol   string
	Period       period.Symbol
	timeserie.TimeSerie[Candlestick]
}

func NewEmptyList(exchangeName, pairSymbol string, period period.Symbol) *List {
	return &List{
		ExchangeName: exchangeName,
		PairSymbol:   pairSymbol,
		Period:       period,
		TimeSerie:    *timeserie.New[Candlestick](),
	}
}

func NewEmptyListFrom(l *List) *List {
	return NewEmptyList(l.ExchangeName, l.PairSymbol, l.Period)
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
	if l.ExchangeName != l2.ExchangeName {
		return ErrExchangeMismatch
	} else if l.PairSymbol != l2.PairSymbol {
		return ErrPairMismatch
	} else if l.Period != l2.Period {
		return ErrPeriodMismatch
	}

	return l.TimeSerie.Merge(l2.TimeSerie, options)
}

func (l List) Extract(start, end time.Time, limit uint) *List {
	el := NewEmptyList(l.ExchangeName, l.PairSymbol, l.Period)
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

func (l *List) HasUncomplete() bool {
	hasUncomplete := false

	_ = l.Loop(func(t time.Time, cs Candlestick) (bool, error) {
		if cs.Uncomplete {
			hasUncomplete = true
			return true, nil
		}
		return false, nil
	})

	return hasUncomplete
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
	txt := fmt.Sprintf("# %s - %s - %s\n", l.ExchangeName, l.PairSymbol, l.Period.String())

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

// AreMissing checks if there is missing candlesticks between two times
// Time order: start < end
func (cl List) AreMissing(end, start time.Time, limit uint) bool {
	if missing := cl.TimeSerie.AreMissing(end, start, cl.Period.Duration(), limit); missing {
		return true
	}

	if cl.HasUncomplete() {
		return true
	}

	return false
}
