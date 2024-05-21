// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=forwardtests.go -destination=mock.gen.go -package client

package client

import (
	"context"

	"github.com/google/uuid"
	client "github.com/lerenn/cryptellation/pkg/client"
	"github.com/lerenn/cryptellation/svc/forwardtests/pkg/forwardtest"
)

type Client interface {
	CreateForwardTest(ctx context.Context, payload forwardtest.NewPayload) (uuid.UUID, error)

	ServiceInfo(ctx context.Context) (client.ServiceInfo, error)
	Close(ctx context.Context)
}
