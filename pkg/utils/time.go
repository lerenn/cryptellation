package utils

import (
	"time"
)

// CountBetweenTimes returns the number of interval between two times
// rounded down to the given interval.
func CountBetweenTimes(t1, t2 time.Time, interval time.Duration) int64 {
	roundedT1 := RoundDownTime(t1, interval)
	roundedT2 := RoundDownTime(t2, interval)

	count := (roundedT2.Unix() - roundedT1.Unix()) / int64(interval/time.Second)
	if count < 0 {
		return -count
	}
	return count
}

// RoundDownTime returns the given time rounded down to the given interval.
func RoundDownTime(t time.Time, interval time.Duration) time.Time {
	diff := t.Unix() % int64(interval/time.Second)
	return time.Unix(t.Unix()-diff, 0)
}

// ElapsedTime returns the time elapsed to execute the given callback
// and the error returned by the callback.
func ElapsedTime(callback func() error) (time.Duration, error) {
	start := time.Now()
	err := callback()
	return time.Since(start), err
}
