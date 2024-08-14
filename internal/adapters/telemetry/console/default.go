package console

import (
	"context"

	"cryptellation/internal/adapters/telemetry"
)

// Fallback to console telemeter if the given telemeter is nil or err is not nil.
// Telemeter will be set to the global telemetry package.
func Fallback(t telemetry.Telemeter, err error) {
	if t == nil || err != nil {
		t = Telemeter{}
		t.Logger(context.Background()).Warning("error with given telemetry, fallback to console telemetry")
	}

	telemetry.Set(t)
}
