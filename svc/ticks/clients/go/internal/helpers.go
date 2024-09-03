package internal

import (
	client "github.com/lerenn/cryptellation/svc/ticks/clients/go"
	"github.com/lerenn/cryptellation/svc/ticks/clients/go/internal/retry"
)

func WrapWithHelpers(client client.Client) client.Client {
	return retry.New(client)
}
