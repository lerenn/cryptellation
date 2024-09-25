package internal

import (
	client "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"
	"github.com/lerenn/cryptellation/svc/candlesticks/clients/go/internal/cache"
	"github.com/lerenn/cryptellation/svc/candlesticks/clients/go/internal/retry"
	"github.com/lerenn/cryptellation/svc/candlesticks/clients/go/internal/splitter"
)

func WrapWithHelpers(client client.Client) client.Client {
	client = retry.New(client)
	client = cache.New(client)
	client = splitter.New(client)
	return client
}
