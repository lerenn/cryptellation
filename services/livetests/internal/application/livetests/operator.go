package livetests

import (
	"context"

	"github.com/digital-feather/cryptellation/services/livetests/internal/domain/livetest"
)

type Operator interface {
	Create(ctx context.Context, req livetest.NewPayload) (id uint, err error)
}
