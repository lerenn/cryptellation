package timeserie

import (
	"errors"
	"time"
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
