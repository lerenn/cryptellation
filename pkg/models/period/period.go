package period

import (
	"errors"
	"fmt"
	"time"

	"github.com/lerenn/cryptellation/v1/pkg/utils"
)

var (
	// ErrInvalidPeriod is returned when the period is invalid.
	ErrInvalidPeriod = errors.New("invalid period")
)

// Symbol is a period symbol.
type Symbol string

const (
	// M1 is a period of 1 minute.
	M1 Symbol = "M1"
	// M3 is a period of 3 minutes.
	M3 Symbol = "M3"
	// M5 is a period of 5 minutes.
	M5 Symbol = "M5"
	// M15 is a period of 15 minutes.
	M15 Symbol = "M15"
	// M30 is a period of 30 minutes.
	M30 Symbol = "M30"
	// H1 is a period of 1 hour.
	H1 Symbol = "H1"
	// H2 is a period of 2 hours.
	H2 Symbol = "H2"
	// H4 is a period of 4 hours.
	H4 Symbol = "H4"
	// H6 is a period of 6 hours.
	H6 Symbol = "H6"
	// H8 is a period of 8 hours.
	H8 Symbol = "H8"
	// H12 is a period of 12 hours.
	H12 Symbol = "H12"
	// D1 is a period of 1 day.
	D1 Symbol = "D1"
	// D3 is a period of 3 days.
	D3 Symbol = "D3"
	// W1 is a period of 1 week.
	W1 Symbol = "W1"
)

// String returns the string representation of the symbol.
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

// Duration converts a symbol into its corresponding duration.
func (s Symbol) Duration() time.Duration {
	return symbolToDuration[s]
}

// FromString converts a string in a symbol and checks if it's valid.
func FromString(symbol string) (Symbol, error) {
	s := Symbol(symbol)
	return s, s.Validate()
}

// FromDuration converts a duration into its corresponding symbol.
func FromDuration(d time.Duration) (Symbol, error) {
	s, ok := durationToSymbol[d]
	if !ok {
		return "", fmt.Errorf("getting symbol from duration (%s): %w", d, ErrInvalidPeriod)
	}

	return s, nil
}

// Validate checks the symbol is a valid period.
func (s Symbol) Validate() error {
	_, ok := symbolToDuration[s]
	if !ok {
		return fmt.Errorf("parsing period from name (%s): %w", s, ErrInvalidPeriod)
	}

	return nil
}

// Symbols returns all the available periods.
func Symbols() []Symbol {
	durations := make([]Symbol, 0, len(symbolToDuration))
	for s := range symbolToDuration {
		durations = append(durations, s)
	}
	return durations
}

// RoundTime rounds the given time to the closest period.
func (s Symbol) RoundTime(t time.Time) time.Time {
	diff := t.Unix() % int64(s.Duration()/time.Second)
	return time.Unix(t.Unix()-diff, 0)
}

// IsAligned checks if the given time is aligned with the period.
func (s Symbol) IsAligned(t time.Time) bool {
	return (t.Unix() % int64(s.Duration()/time.Second)) == 0
}

// FromSeconds returns a period from seconds.
func FromSeconds(i int64) (Symbol, error) {
	for s, p := range symbolToDuration {
		if p == time.Duration(i)*time.Second {
			return s, nil
		}
	}

	return Symbol(""), fmt.Errorf("parsing period from seconds (%d): %w", i, ErrInvalidPeriod)
}

// CountBetweenTimes returns the number of candlesticks between times.
func (s Symbol) CountBetweenTimes(t1, t2 time.Time) int64 {
	return utils.CountBetweenTimes(t1, t2, s.Duration())
}

// UniqueArray returns a unique array of symbols.
func UniqueArray(sym1, sym2 []Symbol) []Symbol {
	return utils.MergeSliceIntoUnique(sym1, sym2)
}

// Opt returns a pointer to the symbol.
func (s Symbol) Opt() *Symbol {
	return &s
}

// RoundInterval takes an interval (represented by start, end) and
// returns the closest time before or equal to 'start' corresponding to the period and the
// closest time after or equal to 'end' corresponding to the period.
//
// Example: if M1 then 1:30 to 2:30 will become 1:00 to 3:00.
func (s Symbol) RoundInterval(start, end *time.Time) (time.Time, time.Time) {
	var nstart, nend time.Time

	defaultDuration := s.Duration() * 500
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

	nstart = s.RoundTime(nstart)
	nend = s.RoundTime(nend)

	return nstart, nend
}
