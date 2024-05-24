// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=port.go -destination=mock.gen.go -package db

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/svc/forwardtests/pkg/forwardtest"
)

type Port interface {
	CreateForwardTest(ctx context.Context, ft forwardtest.ForwardTest) error
	ReadForwardTest(ctx context.Context, id uuid.UUID) (forwardtest.ForwardTest, error)
	ListForwardTests(ctx context.Context, filters ListFilters) ([]forwardtest.ForwardTest, error)
	UpdateForwardTest(ctx context.Context, ft forwardtest.ForwardTest) error
	DeleteForwardTest(ctx context.Context, id uuid.UUID) error
}

type ListFilters struct {
}
