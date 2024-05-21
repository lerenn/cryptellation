package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/svc/forwardtests/pkg/forwardtest"
)

func (f ForwardTests) Create(ctx context.Context, payload forwardtest.NewPayload) (uuid.UUID, error) {
	if err := payload.Validate(); err != nil {
		return uuid.Nil, err
	}

	ft := forwardtest.New(payload)
	return ft.ID, f.db.CreateForwardTest(ctx, ft)
}
