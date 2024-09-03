package cache

import (
	"context"
	"errors"
	"time"

	common "github.com/lerenn/cryptellation/pkg/client"

	client "github.com/lerenn/cryptellation/svc/exchanges/clients/go"
	"github.com/lerenn/cryptellation/svc/exchanges/pkg/exchange"

	"github.com/bluele/gcache"
)

const (
	DefaultMaxSize        = 10000
	DefaultExpirationTime = time.Hour
)

type cache struct {
	client   client.Client
	cache    gcache.Cache
	settings struct {
		maxSize        int
		expirationTime time.Duration
	}
}

func New(client client.Client, options ...option) *cache {
	var c cache

	// Set client and default params
	c.client = client
	c.settings.maxSize = DefaultMaxSize
	c.settings.expirationTime = DefaultExpirationTime

	// Execute options
	for _, option := range options {
		option(&c)
	}

	// Set cache
	c.cache = gcache.New(c.settings.maxSize).LRU().Build()

	return &c
}

func (client *cache) Read(ctx context.Context, names ...string) ([]exchange.Exchange, error) {
	list := make([]exchange.Exchange, 0, len(names))

	missingExchanges := make([]string, 0, len(names))
	for _, name := range names {
		e, err := client.cache.Get(name)
		if errors.Is(err, gcache.KeyNotFoundError) {
			missingExchanges = append(missingExchanges, name)
			continue
		} else if err != nil { // If there is an error
			return nil, err
		}

		// Check if uncomplete
		exch := e.(exchange.Exchange)
		if exch.LastSyncTime.Before(time.Now().Add(-client.settings.expirationTime)) {
			missingExchanges = append(missingExchanges, name)
			continue
		}

		// Add to list
		list = append(list, exch)
	}

	// Return list if none is missing
	if len(missingExchanges) == 0 {
		return list, nil
	}

	// Get missing exchanges
	exchanges, err := client.client.Read(ctx, missingExchanges...)
	if err != nil {
		return nil, err
	}

	// Add to cache
	for _, exch := range exchanges {
		if err := client.cache.Set(exch.Name, exch); err != nil {
			return nil, err
		}
	}

	// Add to list and return
	return append(list, exchanges...), nil
}

func (client *cache) ServiceInfo(ctx context.Context) (common.ServiceInfo, error) {
	return client.client.ServiceInfo(ctx)
}

func (client *cache) Close(ctx context.Context) {
	client.cache.Purge()
	client.client.Close(ctx)
}
