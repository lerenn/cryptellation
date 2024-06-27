package app

import (
	"context"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/pkg/models/account"
	"github.com/lerenn/cryptellation/pkg/models/order"
	"github.com/lerenn/cryptellation/svc/forwardtests/pkg/forwardtest"
)

type ForwardTests interface {
	Create(context.Context, forwardtest.NewPayload) (uuid.UUID, error)
	List(context.Context, ListFilters) ([]forwardtest.ForwardTest, error)
	CreateOrder(ctx context.Context, forwardTest uuid.UUID, order order.Order) error
	GetAccounts(ctx context.Context, forwardTest uuid.UUID) (map[string]account.Account, error)
	GetStatus(ctx context.Context, forwardTest uuid.UUID) (forwardtest.Status, error)
}

type ListFilters struct {
}
