package period

import (
	"errors"
	"fmt"
	"time"

	"cryptellation/pkg/utils"
)

var (
	ErrInvalidPeriod = errors.New("invalid period")
)

type Symbol string

const (
	M1  Symbol = "M1"
	M3  Symbol = "M3"
	M5  Symbol = "M5"
	M15 Symbol = "M15"
	M30 Symbol = "M30"
	H1  Symbol = "H1"
	H2  Symbol = "H2"
	H4  Symbol = "H4"
	H6  Symbol = "H6"
	H8  Symbol = "H8"
	H12 Symbol = "H12"
	D1  Symbol = "D1"
	D3  Symbol = "D3"
	W1  Symbol = "W1"
)

func (s Symbol) String() string {
	return string(s)
}

var (
	symbolToDuration = map[Symbol]time.Duration{
		M1:  time.Minute,
		M3:  3 * time.Minute,
		M5:  5 * time.Minute,
		M15: 15 * time.Minute,
		M30: 30 * time.Minute,
		H1:  time.Hour,
		H2:  2 * time.Hour,
		H4:  4 * time.Hour,
		H6:  6 * time.Hour,
		H8:  8 * time.Hour,
		H12: 12 * time.Hour,
		D1:  24 * time.Hour,
		D3:  3 * 24 * time.Hour,
		W1:  7 * 24 * time.Hour,
	}

	durationToSymbol = map[time.Duration]Symbol{
		time.Minute:        M1,
		3 * time.Minute:    M3,
		5 * time.Minute:    M5,
		15 * time.Minute:   M15,
		30 * time.Minute:   M30,
		time.Hour:          H1,
		2 * time.Hour:      H2,
		4 * time.Hour:      H4,
		6 * time.Hour:      H6,
		8 * time.Hour:      H8,
		12 * time.Hour:     H12,
		24 * time.Hour:     D1,
		3 * 24 * time.Hour: D3,
		7 * 24 * time.Hour: W1,
	}
)

func (s Symbol) Duration() time.Duration {
	return symbolToDuration[s]
}

func FromString(symbol string) (Symbol, error) {
	s := Symbol(symbol)
	return s, s.Validate()
}

func FromDuration(d time.Duration) (Symbol, error) {
	s, ok := durationToSymbol[d]
	if !ok {
		return "", fmt.Errorf("getting symbol from duration (%s): %w", d, ErrInvalidPeriod)
	}

	return s, nil
}

func (s Symbol) Validate() error {
	_, ok := symbolToDuration[s]
	if !ok {
		return fmt.Errorf("parsing period from name (%s): %w", s, ErrInvalidPeriod)
	}

	return nil
}

func Symbols() []Symbol {
	durations := make([]Symbol, 0, len(symbolToDuration))
	for s := range symbolToDuration {
		durations = append(durations, s)
	}
	return durations
}

func (s Symbol) RoundTime(t time.Time) time.Time {
	diff := t.Unix() % int64(s.Duration()/time.Second)
	return time.Unix(t.Unix()-diff, 0)
}

func (s Symbol) IsAligned(t time.Time) bool {
	return (t.Unix() % int64(s.Duration()/time.Second)) == 0
}

func FromSeconds(i int64) (Symbol, error) {
	for s, p := range symbolToDuration {
		if p == time.Duration(i)*time.Second {
			return s, nil
		}
	}

	return Symbol(""), fmt.Errorf("parsing period from seconds (%d): %w", i, ErrInvalidPeriod)
}

// CountBetweenTimes returns the number of candlesticks between times
func (s Symbol) CountBetweenTimes(t1, t2 time.Time) int64 {
	return utils.CountBetweenTimes(t1, t2, s.Duration())
}

func UniqueArray(sym1, sym2 []Symbol) []Symbol {
	tmp := make([]Symbol, len(sym1))
	copy(tmp, sym1)

	for _, s2 := range sym2 {
		present := false
		for _, s1 := range sym1 {
			if s1 == s2 {
				present = true
				break
			}
		}

		if !present {
			tmp = append(tmp, s2)
		}
	}

	return tmp
}

func (s Symbol) Opt() *Symbol {
	return &s
}

// RoundInterval takes an interval (represented by start, end) and
// returns the closest time before or equal to 'start' corresponding to the period and the
// closest time after or equal to 'end' corresponding to the period.
//
// Example: if M1 then 1:30 to 2:30 will become 1:00 to 3:00.
func (per Symbol) RoundInterval(start, end *time.Time) (time.Time, time.Time) {
	var nstart, nend time.Time

	defaultDuration := per.Duration() * 500
	if end == nil {
		if start == nil {
			nend = time.Now()
		} else {
			nend = start.Add(defaultDuration)
		}
	} else {
		nend = *end
	}

	if start == nil {
		nstart = nend.Add(-defaultDuration)
	} else {
		nstart = *start
	}

	nstart = per.RoundTime(nstart)
	nend = per.RoundTime(nend)

	return nstart, nend
}
