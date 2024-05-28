package timeserie

import (
	"fmt"
	"slices"
	"time"
)

var (
	ErrTimeRangeInvalid       = fmt.Errorf("time range invalid")
	ErrTimeRangeStartAfterEnd = fmt.Errorf("%w: start after end", ErrTimeRangeInvalid)
)

// TimeRange is a structure representing the slot of time, from start to end
type TimeRange struct {
	Start, End time.Time
}

func (tr TimeRange) Validate() error {
	if tr.Start.After(tr.End) {
		s := tr.Start.Format(time.RFC3339)
		e := tr.End.Format(time.RFC3339)
		return fmt.Errorf("%w: start=%s, end=%s", ErrTimeRangeStartAfterEnd, s, e)
	}

	return nil
}

func TimeRangesToString(tr []TimeRange) string {
	var str string
	for _, t := range tr {
		str = fmt.Sprintf("[%s - %s]", t.Start.Format(time.RFC3339), t.End.Format(time.RFC3339))
	}
	return str
}

func MergeTimeRanges(tr1, tr2 []TimeRange) ([]TimeRange, error) {
	// Check time range 1
	for i, tr := range tr1 {
		if err := tr.Validate(); err != nil {
			return nil, fmt.Errorf("error on tr1 %d: %w", i, err)
		}
	}

	// Check time range 2
	for i, tr := range tr1 {
		if err := tr.Validate(); err != nil {
			return nil, fmt.Errorf("error on tr1 %d: %w", i, err)
		}
	}

	// Merge time ranges ordered
	return mergeTimeRangesWithoutOrdering(
		OrderTimeRanges(tr1),
		OrderTimeRanges(tr2),
	), nil
}

func mergeTimeRangesWithoutOrdering(tr1, tr2 []TimeRange) []TimeRange {
	if len(tr1) == 0 {
		return tr2
	} else if len(tr2) == 0 {
		return tr1
	}

	resulting := make([]TimeRange, 0, len(tr1)+len(tr2))
	tr1Current, tr2Current := 0, 0

	var lastStart time.Time
	var cursor struct {
		TR1, TR2 bool
	}
	if tr1[0].Start.Before(tr2[0].Start) {
		lastStart = tr1[0].Start
		cursor.TR1 = true
	} else {
		lastStart = tr2[0].Start
		cursor.TR2 = true
	}

	for {
		// Both are finished
		if tr1Current >= len(tr1) && tr2Current >= len(tr2) {
			return resulting
		}

		// One of them is finished
		if tr1Current >= len(tr1) {
			if cursor.TR2 {
				resulting = append(resulting, TimeRange{
					Start: lastStart,
					End:   tr2[tr2Current].End,
				})
				tr2Current++
			}
			resulting = append(resulting, tr2[tr2Current:]...)
		} else if tr2Current >= len(tr2) {
			if cursor.TR1 {
				resulting = append(resulting, TimeRange{
					Start: lastStart,
					End:   tr1[tr1Current].End,
				})
				tr1Current++
			}
			resulting = append(resulting, tr1[tr1Current:]...)
		}
		if tr1Current >= len(tr1) || tr2Current >= len(tr2) {
			return resulting
		}

		if cursor.TR1 && cursor.TR2 {
			if tr1[tr1Current].End.After(tr2[tr2Current].End) {
				cursor.TR2 = false
				tr2Current++
			} else if tr1[tr1Current].End.Before(tr2[tr2Current].End) {
				cursor.TR1 = false
				tr1Current++
			} else {
				cursor.TR1 = false
				cursor.TR2 = false
				resulting = append(resulting, TimeRange{
					Start: lastStart,
					End:   tr1[tr1Current].End,
				})
				tr1Current++
				tr2Current++
			}
		} else if cursor.TR1 {
			if tr1[tr1Current].End.After(tr2[tr2Current].Start) {
				cursor.TR2 = true
			} else if tr1[tr1Current].End.Before(tr2[tr2Current].Start) {
				cursor.TR1 = false
				resulting = append(resulting, TimeRange{
					Start: lastStart,
					End:   tr1[tr1Current].End,
				})
				tr1Current++
			} else {
				cursor.TR1 = false
				tr1Current++
				cursor.TR2 = true
			}
		} else if cursor.TR2 {
			if tr1[tr1Current].Start.After(tr2[tr2Current].End) {
				cursor.TR2 = false
				resulting = append(resulting, TimeRange{
					Start: lastStart,
					End:   tr2[tr2Current].End,
				})
				tr2Current++
			} else if tr1[tr1Current].Start.Before(tr2[tr2Current].End) {
				cursor.TR1 = true
			} else {
				cursor.TR1 = true
				cursor.TR2 = false
				tr2Current++
			}
		} else {
			if tr1[tr1Current].Start.After(tr2[tr2Current].Start) {
				lastStart = tr2[tr2Current].Start
				cursor.TR2 = true
			} else if tr1[tr1Current].Start.Before(tr2[tr2Current].Start) {
				lastStart = tr1[tr1Current].Start
				cursor.TR1 = true
			} else {
				lastStart = tr1[tr1Current].Start
				cursor.TR1 = true
				cursor.TR2 = true
			}
		}
	}
}

func OrderTimeRanges(tr []TimeRange) []TimeRange {
	slices.SortFunc(tr, func(a, b TimeRange) int {
		if a.Start.Before(b.Start) {
			return -1
		} else if a.Start.Before(b.Start) {
			return 1
		} else {
			return 0
		}
	})

	return mergeTimeRangesWithoutOrdering(tr, tr)
}

func TimeRangesFromMissingTimes(times []time.Time) []TimeRange {
	if len(times) == 0 {
		return []TimeRange{}
	}

	tr := make([]TimeRange, 0, len(times))
	current := TimeRange{
		Start: times[0],
		End:   times[0],
	}

	for _, t := range times[1:] {
		if current.End.Add(1).Equal(t) {
			current.End = t
			continue
		}

		tr = append(tr, current)
		current.Start, current.End = t, t
	}

	return append(tr, current)
}
