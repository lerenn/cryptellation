package console

import (
	"github.com/lerenn/cryptellation/v1/pkg/telemetry"
)

// Fallback to console telemeter if the given telemeter is nil or err is not nil.
// Telemeter will be set to the global telemetry package.
func Fallback(t telemetry.Telemeter, err error) {
	if t == nil || err != nil {
		t = Telemeter{}
	}

	telemetry.Set(t)
}
