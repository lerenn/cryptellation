package internal

import (
	client "github.com/lerenn/cryptellation/svc/exchanges/clients/go"
	"github.com/lerenn/cryptellation/svc/exchanges/clients/go/internal/cache"
	"github.com/lerenn/cryptellation/svc/exchanges/clients/go/internal/retry"
)

func WrapWithHelpers(client client.Client) client.Client {
	client = retry.New(client)
	return cache.New(client)
}
