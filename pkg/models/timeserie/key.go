package timeserie

import "time"

type key int64

func newKey(t time.Time) key {
	return key(t.UnixNano())
}

func (k key) ToTime() time.Time {
	return time.Unix(0, int64(k))
}
