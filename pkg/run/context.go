package run

import (
	"time"

	"github.com/google/uuid"
)

// Context is the context of a run that contains several information about the run.
type Context struct {
	ID        uuid.UUID
	Mode      Mode
	Now       time.Time
	TaskQueue string
}
