package client

import (
	"context"
	"errors"
	"time"

	"github.com/bluele/gcache"
	"github.com/lerenn/cryptellation/pkg/client"
	"github.com/lerenn/cryptellation/svc/exchanges/pkg/exchange"
)

var _ Client = (*CachedClient)(nil)

type CachedClient struct {
	controller Client
	cache      gcache.Cache
	parameters CacheParameters
}

type CacheParameters struct {
	MaxSize        int
	ExpirationTime time.Duration
}

const (
	DefaultMaxSize        = 10000
	DefaultExpirationTime = time.Hour
)

func DefaultCacheParameters() CacheParameters {
	return CacheParameters{
		MaxSize:        DefaultMaxSize,
		ExpirationTime: DefaultExpirationTime,
	}
}

func NewCachedClient(controller Client, params CacheParameters) *CachedClient {
	return &CachedClient{
		controller: controller,
		cache:      gcache.New(params.MaxSize).LRU().Build(),
		parameters: params,
	}
}

func (client *CachedClient) Read(ctx context.Context, names ...string) ([]exchange.Exchange, error) {
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
		if exch.LastSyncTime.Before(time.Now().Add(-client.parameters.ExpirationTime)) {
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
	exchanges, err := client.controller.Read(ctx, missingExchanges...)
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

func (client *CachedClient) ServiceInfo(ctx context.Context) (client.ServiceInfo, error) {
	return client.controller.ServiceInfo(ctx)
}

func (client *CachedClient) Close(ctx context.Context) {
	client.cache.Purge()
	client.controller.Close(ctx)
}
