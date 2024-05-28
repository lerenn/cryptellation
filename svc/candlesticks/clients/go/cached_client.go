package client

import (
	"context"
	"errors"
	"time"

	"github.com/bluele/gcache"
	client "github.com/lerenn/cryptellation/pkg/client"
	"github.com/lerenn/cryptellation/pkg/models/timeserie"
	"github.com/lerenn/cryptellation/pkg/utils"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"
)

var _ Client = (*CachedClient)(nil)

type cacheKey struct {
	Exchange  string
	Pair      string
	Period    period.Symbol
	Timestamp int64
}

type CachedClient struct {
	controller Client
	cache      gcache.Cache
	parameters CacheParameters
}

type CacheParameters struct {
	MaxSize              int
	PreLoadingAfterSize  int
	PreLoadingBeforeSize int
}

const (
	DefaultMaxSize              = 10000
	DefaultPreLoadingAfterSize  = 200
	DefaultPreLoadingBeforeSize = 0
)

func DefaultCacheParameters() CacheParameters {
	return CacheParameters{
		MaxSize:              DefaultMaxSize,
		PreLoadingAfterSize:  DefaultPreLoadingAfterSize,
		PreLoadingBeforeSize: DefaultPreLoadingBeforeSize,
	}
}

func NewCachedClient(controller Client, params CacheParameters) *CachedClient {
	return &CachedClient{
		controller: controller,
		cache:      gcache.New(params.MaxSize).LRU().Build(),
		parameters: params,
	}
}

type ReadCandlesticksPayload struct {
	Exchange string
	Pair     string
	Period   period.Symbol
	Start    *time.Time
	End      *time.Time
	Limit    uint
}

func (client *CachedClient) Read(ctx context.Context, payload ReadCandlesticksPayload) (*candlestick.List, error) {
	list := candlestick.NewList(payload.Exchange, payload.Pair, payload.Period)

	if payload.Start == nil {
		return nil, errors.New("payload.Start is required")
	} else if payload.End == nil {
		return nil, errors.New("payload.End is required")
	}

	// Get missing times
	missingTimes := make([]time.Time, 0, payload.Period.CountBetweenTimes(*payload.Start, *payload.End))
	for current := *payload.Start; !current.After(*payload.End); current = current.Add(payload.Period.Duration()) {
		key := cacheKey{
			Exchange:  payload.Exchange,
			Pair:      payload.Pair,
			Period:    payload.Period,
			Timestamp: current.Unix(),
		}
		c, err := client.cache.Get(key)
		if errors.Is(err, gcache.KeyNotFoundError) {
			missingTimes = append(missingTimes, current)
		} else if err != nil {
			return nil, err
		} else if err := list.Set(current, c.(candlestick.Candlestick)); err != nil {
			return nil, err
		}
	}

	// Change to time ranges and return if none
	tr := timeserie.TimeRangesFromMissingTimes(missingTimes)
	if len(tr) == 0 {
		return list, nil
	}

	// Generate new payload with extended time ranges
	newPayload := payload
	newPayload.Start = utils.ToReference(tr[0].Start.Add(-payload.Period.Duration() * time.Duration(client.parameters.PreLoadingBeforeSize)))
	newPayload.End = utils.ToReference(tr[len(tr)-1].End.Add(payload.Period.Duration() * time.Duration(client.parameters.PreLoadingAfterSize)))

	// Get missing times
	missing, err := client.controller.Read(ctx, newPayload)
	if err != nil {
		return nil, err
	}

	// Add missing times to cache
	if err := missing.Loop(func(t time.Time, c candlestick.Candlestick) (bool, error) {
		key := cacheKey{
			Exchange:  payload.Exchange,
			Pair:      payload.Pair,
			Period:    payload.Period,
			Timestamp: t.Unix(),
		}
		return false, client.cache.Set(key, c)
	}); err != nil {
		return nil, err
	}

	// Merge missing times
	if err := list.Merge(missing, &timeserie.MergeOptions{}); err != nil {
		return nil, err
	}

	// Exctract only requested
	return list.Extract(*payload.Start, *payload.End, 0), nil
}

func (client *CachedClient) ServiceInfo(ctx context.Context) (client.ServiceInfo, error) {
	return client.controller.ServiceInfo(ctx)
}

func (client *CachedClient) Close(ctx context.Context) {
	client.controller.Close(ctx)
}
