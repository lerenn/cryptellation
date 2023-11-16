package telemetry

import "context"

type Telemetry interface {
	Close(ctx context.Context)
}
