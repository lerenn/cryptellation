package temporal

import (
	context "context"
	"time"

	"go.temporal.io/sdk/activity"
)

// AsyncActivityHeartbeat will start a goroutine that will send a heartbeat every
// timeout until the context is done.
func AsyncActivityHeartbeat(ctx context.Context, timeout time.Duration) {
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(timeout):
				activity.RecordHeartbeat(ctx, nil)
			}
		}
	}(ctx)
}
