package internal

import (
	client "github.com/lerenn/cryptellation/svc/forwardtests/clients/go"
	"github.com/lerenn/cryptellation/svc/forwardtests/clients/go/internal/retry"
)

func WrapWithHelpers(client client.Client) client.Client {
	return retry.New(client)
}
