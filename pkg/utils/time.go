package utils

import (
	"time"
)

// CountBetweenTimes will return the number of interval between two times
// rounded down to the given interval
func CountBetweenTimes(t1, t2 time.Time, interval time.Duration) int64 {
	roundedT1 := RoundDownTime(t1, interval)
	roundedT2 := RoundDownTime(t2, interval)
	count := (roundedT2.Unix() - roundedT1.Unix()) / int64(interval/time.Second)
	if count < 0 {
		return -count
	} else {
		return count
	}
}

func RoundDownTime(t time.Time, interval time.Duration) time.Time {
	diff := t.Unix() % int64(interval/time.Second)
	return time.Unix(t.Unix()-diff, 0)
}

func ElapsedTime(callback func() error) (time.Duration, error) {
	start := time.Now()
	err := callback()
	return time.Since(start), err
}
