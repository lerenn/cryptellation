package timeserie

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/lerenn/cryptellation/v1/pkg/utils"
)

var (
	// ErrTimeStampAlreadyExists is returned when trying to set a timestamp that already exists.
	ErrTimeStampAlreadyExists = errors.New("timestamp-already-exists")
)

// TimeSerie is a time serie data structure that can be used to store data over time.
type TimeSerie[T any] struct {
	data        map[key]T
	orderedKeys []key
}

// New creates a new TimeSerie.
func New[T any]() *TimeSerie[T] {
	return &TimeSerie[T]{
		data:        make(map[key]T),
		orderedKeys: make([]key, 0),
	}
}

// MarshalJSON marshals the TimeSerie into a JSON object.
func (ts *TimeSerie[T]) MarshalJSON() ([]byte, error) {
	var data = make(map[string]T)

	if err := ts.Loop(func(t time.Time, d T) (bool, error) {
		data[t.Format(time.RFC3339)] = d
		return false, nil
	}); err != nil {
		return nil, err
	}

	return json.Marshal(data)
}

// UnmarshalJSON unmarshals the TimeSerie from a JSON object.
func (ts *TimeSerie[T]) UnmarshalJSON(data []byte) error {
	var st = make(map[string]T)

	// Unmarshal the data into a map[string]T
	if err := json.Unmarshal(data, &st); err != nil {
		return err
	}

	// Initialize the TimeSerie
	*ts = *New[T]()

	// Loop over the map and set the data
	for k, v := range st {
		t, err := time.Parse(time.RFC3339, k)
		if err != nil {
			return err
		}

		ts.Set(t, v)
	}

	return nil
}

// Set sets a value at a specific time.
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

	switch {
	case previous == -1:
		ts.orderedKeys = append([]key{k}, ts.orderedKeys...)
	case previous == len(ts.orderedKeys)-1:
		ts.orderedKeys = append(ts.orderedKeys, k)
	default:
		index := previous + 1
		ts.orderedKeys = append(ts.orderedKeys[:index+1], ts.orderedKeys[index:]...)
		ts.orderedKeys[index] = k
	}
}

// Get gets a value at a specific time. The second return value is false if the
// value for this specific time does not exist.
func (ts *TimeSerie[T]) Get(t time.Time) (T, bool) {
	data, exists := ts.data[newKey(t)]
	return data, exists
}

// Len returns the length of the TimeSerie.
func (ts *TimeSerie[T]) Len() int {
	return len(ts.orderedKeys)
}

// MergeOptions are options for the Merge function.
type MergeOptions struct {
	// ErrorOnCollision returns an error if there is a collision between the two TimeSeries.
	ErrorOnCollision bool
}

// Merge merges two TimeSeries together.
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

// Delete deletes a value at a specific time.
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

		switch {
		case index == 0:
			ts.orderedKeys = ts.orderedKeys[1:]
		case index == len(ts.orderedKeys)-1:
			ts.orderedKeys = ts.orderedKeys[:len(ts.orderedKeys)-1]
		default:
			ts.orderedKeys = append(ts.orderedKeys[:index], ts.orderedKeys[index+1:]...)
		}
	}
}

// Loop loops over every elements of the TimeSerie in time order and calls the callback function.
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

// First returns the first element of the TimeSerie.
func (ts TimeSerie[T]) First() (time.Time, T, bool) {
	if ts.Len() == 0 {
		return time.Unix(0, 0), *new(T), false
	}

	t := ts.orderedKeys[0]
	return t.ToTime(), ts.data[t], true
}

// Last returns the last element of the TimeSerie.
func (ts TimeSerie[T]) Last() (time.Time, T, bool) {
	if ts.Len() == 0 {
		return time.Unix(0, 0), *new(T), false
	}

	t := ts.orderedKeys[ts.Len()-1]
	return t.ToTime(), ts.data[t], true
}

// Extract extracts a TimeSerie between two times.
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

// FirstN returns the first N elements of the TimeSerie.
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

// AreMissing checks if there is missing candlesticks between two times.
func (ts TimeSerie[T]) AreMissing(start, end time.Time, interval time.Duration, limit uint) bool {
	expectedCount := int(utils.CountBetweenTimes(start, end, interval)) + 1
	qty := ts.Len()

	if qty < expectedCount && (limit == 0 || uint(qty) < limit) {
		return true
	}

	return false
}

// GetMissingTimes returns an array of missing time in the timeserie.
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

// GetMissingRanges returns an array of missing time ranges in the timeserie.
func (ts TimeSerie[T]) GetMissingRanges(start, end time.Time, interval time.Duration, limit uint) []TimeRange {
	missing := ts.GetMissingTimes(start, end, interval, limit)
	return TimeRangesFromMissingTimes(interval, missing)
}
