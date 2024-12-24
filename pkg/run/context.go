package run

import "github.com/google/uuid"

// Context is the context of a run that contains several information about the run.
type Context struct {
	ID        uuid.UUID
	Mode      Mode
	TaskQueue string
}
