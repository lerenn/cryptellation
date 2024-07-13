package client

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

func TestRetryableSuite(t *testing.T) {
	suite.Run(t, new(RetryableSuite))
}

type RetryableSuite struct {
	suite.Suite
}

func (suite *RetryableSuite) TestExecWithExceedingMaxRetries() {
	// Given a Retryable with 3 max retries and a timeout of 100ms
	retryable := Retryable{
		MaxRetries: 3,
		Timeout:    time.Millisecond * 10,
	}

	// When calling the Retry method with a function that always returns an error
	count := 0
	err := retryable.Exec(context.Background(), func(ctx context.Context) error {
		count++
		<-ctx.Done()
		return ctx.Err()
	})

	// Then the error should be the one returned by the function
	suite.ErrorIs(err, ErrMaxRetriesReached)
	suite.Equal(3, count)
}

func (suite *RetryableSuite) TestExecWithOneRetry() {
	// Given a Retryable with 3 max retries and a timeout of 100ms
	retryable := Retryable{
		MaxRetries: 3,
		Timeout:    time.Millisecond * 10,
	}

	// When calling the Retry method with a function that returns an error once
	count := 0
	err := retryable.Exec(context.Background(), func(ctx context.Context) error {
		count++
		if count >= 1 {
			return nil
		}
		<-ctx.Done()
		return ctx.Err()
	})

	// Then the error should be nil
	suite.NoError(err)
	suite.Equal(1, count)
}
