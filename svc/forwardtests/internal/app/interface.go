package app

import (
	"context"

	"cryptellation/pkg/models/account"
	"cryptellation/pkg/models/order"

	"cryptellation/svc/forwardtests/pkg/forwardtest"

	"github.com/google/uuid"
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
