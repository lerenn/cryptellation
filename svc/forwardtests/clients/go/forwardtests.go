// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=forwardtests.go -destination=mock.gen.go -package client

package client

import (
	"context"

	"cryptellation/pkg/client"
	"cryptellation/pkg/models/account"

	"cryptellation/svc/forwardtests/pkg/forwardtest"

	"github.com/google/uuid"
)

type Client interface {
	CreateForwardTest(ctx context.Context, payload forwardtest.NewPayload) (uuid.UUID, error)
	ListForwardTests(ctx context.Context) ([]forwardtest.ForwardTest, error)
	CreateOrder(ctx context.Context, payload client.OrderCreationPayload) error
	GetAccounts(ctx context.Context, forwardTestID uuid.UUID) (map[string]account.Account, error)
	GetStatus(ctx context.Context, forwardTestID uuid.UUID) (forwardtest.Status, error)

	ServiceInfo(ctx context.Context) (client.ServiceInfo, error)
	Close(ctx context.Context)
}
