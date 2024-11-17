package client

import (
	"context"

	"github.com/lerenn/cryptellation/v1/api"
)

// Client is the client interface.
type Client interface {
	Info(ctx context.Context) (api.ServiceInfoWorkflowResult, error)
	Close(ctx context.Context)
}
