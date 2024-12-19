package run

import "github.com/google/uuid"

type Context struct {
	ID        uuid.UUID
	Mode      Mode
	TaskQueue string
}
