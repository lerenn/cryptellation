package timeserie

import (
	"errors"
	"time"

	"github.com/lerenn/cryptellation/pkg/utils"
)

var (
	ErrTimeStampAlreadyExists error = errors.New("timestamp-already-exists")
)

type TimeSerie[T any] struct {
	data        map[key]T
	orderedKeys []key
}

func New[T any]() *TimeSerie[T] {
	return &TimeSerie[T]{
		data:        make(map[key]T),
		orderedKeys: make([]key, 0),
	}
}

func (ts *TimeSerie[T]) Set(t time.Time, d T) *TimeSerie[T] {
	k := newKey(t)

	ts.addKey(k)
	ts.data[k] = d

	return ts
}

func (ts *TimeSerie[T]) addKey(k key) {
	if _, exists := ts.data[k]; exists {
		return
	}

	previous := -1
	for pos, v := range ts.orderedKeys {
		if v > k {
			break
		}
		previous = pos
	}

	if previous == -1 {
		ts.orderedKeys = append([]key{k}, ts.orderedKeys...)
	} else if previous == len(ts.orderedKeys)-1 {
		ts.orderedKeys = append(ts.orderedKeys, k)
	} else {
		index := previous + 1
		ts.orderedKeys = append(ts.orderedKeys[:index+1], ts.orderedKeys[index:]...)
		ts.orderedKeys[index] = k
	}
}

func (ts *TimeSerie[T]) Get(t time.Time) (T, bool) {
	data, exists := ts.data[newKey(t)]
	return data, exists
}

func (ts *TimeSerie[T]) Len() int {
	return len(ts.orderedKeys)
}

type MergeOptions struct {
	ErrorOnCollision bool
}

func (ts *TimeSerie[T]) Merge(ts2 TimeSerie[T], options *MergeOptions) error {
	if options != nil && options.ErrorOnCollision {
		for pos := range ts2.data {
			if _, exists := ts.Get(pos.ToTime()); exists {
				return ErrTimeStampAlreadyExists
			}
		}
	}

	for pos, v := range ts2.data {
		ts.Set(pos.ToTime(), v)
	}

	return nil
}

func (ts *TimeSerie[T]) Delete(t ...time.Time) {
	for _, ti := range t {
		k := newKey(ti)

		if _, exists := ts.data[k]; !exists {
			return
		}

		delete(ts.data, k)
		index := -1
		for pos, v := range ts.orderedKeys {
			index = pos
			if v == k {
				break
			}
		}

		if index == 0 {
			ts.orderedKeys = ts.orderedKeys[1:]
		} else if index == len(ts.orderedKeys)-1 {
			ts.orderedKeys = ts.orderedKeys[:len(ts.orderedKeys)-1]
		} else {
			ts.orderedKeys = append(ts.orderedKeys[:index], ts.orderedKeys[index+1:]...)
		}
	}
}

func (ts *TimeSerie[T]) Loop(callback func(time.Time, T) (bool, error)) error {
	for _, t := range ts.orderedKeys {
		if shouldBreak, err := callback(t.ToTime(), ts.data[t]); err != nil {
			return err
		} else if shouldBreak {
			break
		}
	}

	return nil
}

func (ts TimeSerie[T]) First() (time.Time, T, bool) {
	if ts.Len() == 0 {
		return time.Unix(0, 0), *new(T), false
	}

	t := ts.orderedKeys[0]
	return t.ToTime(), ts.data[t], true
}

func (ts TimeSerie[T]) Last() (time.Time, T, bool) {
	if ts.Len() == 0 {
		return time.Unix(0, 0), *new(T), false
	}

	t := ts.orderedKeys[ts.Len()-1]
	return t.ToTime(), ts.data[t], true
}

func (ts TimeSerie[T]) Extract(start, end time.Time, limit int) *TimeSerie[T] {
	ets := New[T]()

	count := 0
	_ = ts.Loop(func(t time.Time, obj T) (bool, error) {
		if t.Before(start) {
			return false, nil
		}

		if t.After(end) {
			return true, nil
		}

		ets.Set(t, obj)

		count++
		if limit > 0 && count >= limit {
			return true, nil
		}

		return false, nil
	})

	return ets
}

func (ts TimeSerie[T]) FirstN(limit uint) *TimeSerie[T] {
	ets := New[T]()

	if limit == 0 {
		return ets
	}

	var count uint
	_ = ts.Loop(func(t time.Time, obj T) (bool, error) {
		ets.Set(t, obj)
		count++
		return count >= limit, nil
	})

	return ets
}

// AreMissing checks if there is missing candlesticks between two times
func (ts TimeSerie[T]) AreMissing(start, end time.Time, interval time.Duration, limit uint) bool {
	expectedCount := int(utils.CountBetweenTimes(start, end, interval)) + 1
	qty := ts.Len()

	if qty < expectedCount && (limit == 0 || uint(qty) < limit) {
		return true
	}

	return false
}

// GetMissingTimes returns an array of missing time in the timeserie
func (ts TimeSerie[T]) GetMissingTimes(start, end time.Time, interval time.Duration, limit uint) []time.Time {
	expectedCount := int(utils.CountBetweenTimes(start, end, interval)) + 1
	list := make([]time.Time, 0, expectedCount)
	for current, count := start, uint(0); !current.After(end); current, count = current.Add(interval), count+1 {
		_, exists := ts.Get(current)
		if !exists {
			list = append(list, current)
		}

		if limit > 0 && count >= limit-1 {
			break
		}
	}
	return list
}

func (ts TimeSerie[T]) GetMissingRanges(start, end time.Time, interval time.Duration, limit uint) []TimeRange {
	missing := ts.GetMissingTimes(start, end, interval, limit)
	if len(missing) == 0 {
		return []TimeRange{}
	}

	tr := make([]TimeRange, 0, len(missing))
	current := TimeRange{
		Start: missing[0],
		End:   missing[0],
	}

	for _, t := range missing[1:] {
		if current.End.Add(interval).Equal(t) {
			current.End = t
			continue
		}

		tr = append(tr, current)
		current.Start, current.End = t, t
	}

	return append(tr, current)
}
