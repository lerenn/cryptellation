package app

import (
	"context"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/svc/forwardtests/pkg/forwardtest"
)

type ForwardTests interface {
	Create(context.Context, forwardtest.NewPayload) (uuid.UUID, error)
}
