package retry

import "time"

type option func(r *retry)

func WithMaxRetry(maxRetry int) option {
	return func(r *retry) {
		r.MaxRetries = maxRetry
	}
}

func WithTimeout(timeout time.Duration) option {
	return func(r *retry) {
		r.Timeout = timeout
	}
}
