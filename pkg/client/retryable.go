package client

import (
	"context"
	"errors"
	"time"
)

var (
	ErrMaxRetriesReached = errors.New("max retries reached")
)

type Retryable struct {
	MaxRetries int
	Timeout    time.Duration
}

var DefaultRetryable = Retryable{
	MaxRetries: 3,
	Timeout:    3 * time.Second,
}

func (r Retryable) Exec(ctx context.Context, fn func(ctx context.Context) error) error {
	for i := 0; i < r.MaxRetries; i++ {
		// Create a context with a deadline
		deadlinedContext, cancel := context.WithDeadline(ctx, time.Now().Add(r.Timeout))

		// Execute and cancel the context
		err := fn(deadlinedContext)
		cancel()

		// If there is no error, everything went well
		if err == nil {
			return nil
		}

		// If the error is not a deadline exceeded, then it's a true error
		if !errors.Is(err, context.DeadlineExceeded) {
			return err
		}
	}

	return ErrMaxRetriesReached
}
