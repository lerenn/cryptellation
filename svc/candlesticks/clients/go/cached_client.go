package client

import (
	"context"
	"errors"
	"time"

	"github.com/bluele/gcache"
	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
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
	params     CacheParameters
}

type CacheParameters struct {
	MaxSize                       int
	PreLoadingSize                int
	PreemptiveAsyncLoadingEnabled bool
}

const (
	DefaultMaxSize              = 100000
	DefaultPreLoadingAfterSize  = 200
	DefaultPreLoadingBeforeSize = 0
)

func DefaultCacheParameters() CacheParameters {
	return CacheParameters{
		MaxSize:                       DefaultMaxSize,
		PreLoadingSize:                DefaultPreLoadingAfterSize,
		PreemptiveAsyncLoadingEnabled: true,
	}
}

func NewCachedClient(controller Client, params CacheParameters) *CachedClient {
	return &CachedClient{
		controller: controller,
		cache:      gcache.New(params.MaxSize).LRU().Build(),
		params:     params,
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
	// Check required fields for caching
	if payload.End == nil {
		payload.End = utils.ToReference(time.Now())
	}
	if payload.Start == nil {
		start := payload.End.Add(-payload.Period.Duration() * time.Duration(client.params.PreLoadingSize))
		payload.Start = &start
	}

	// Round down payload start and end
	payload.Start = utils.ToReference(payload.Period.RoundTime(*payload.Start))
	payload.End = utils.ToReference(payload.Period.RoundTime(*payload.End))

	// Get present and missing times
	list, tr, err := client.presentAndmissingTimes(payload)
	if err != nil {
		return nil, err
	}

	if len(tr) == 0 {
		telemetry.L(ctx).Debug("No missing times in cache")
		return list, nil
	}
	telemetry.L(ctx).Debugf("Missing times in cache: %v", tr)

	// Download missing times
	missing, err := client.downloadMissing(ctx, payload, tr)
	if err != nil {
		return nil, err
	}

	// Merge missing times
	if err := list.Merge(missing, &timeserie.MergeOptions{}); err != nil {
		return nil, err
	}

	// Exctract only requested
	extract := list.Extract(*payload.Start, *payload.End, payload.Limit)
	telemetry.L(ctx).Debugf("Returning %d candlesticks from %s to %s", extract.Len(), *payload.Start, *payload.End)

	if client.params.PreemptiveAsyncLoadingEnabled && missing.Len() == 0 {
		go client.preemptiveAsyncLoading(ctx, payload)
	}

	return extract, nil
}

func (client *CachedClient) presentAndmissingTimes(payload ReadCandlesticksPayload) (present *candlestick.List, missing []timeserie.TimeRange, err error) {
	list := candlestick.NewList(payload.Exchange, payload.Pair, payload.Period)

	// Generate missing times slice
	max := payload.Period.CountBetweenTimes(*payload.Start, *payload.End)
	if payload.Limit != 0 && int64(payload.Limit) < max {
		max = int64(payload.Limit)
	}
	missingTimes := make([]time.Time, 0, max)

	// Get missing times
	for count, current := uint(0), *payload.Start; !current.After(*payload.End) && (payload.Limit == 0 || count < payload.Limit); count, current = count+1, current.Add(payload.Period.Duration()) {
		key := cacheKey{
			Exchange:  payload.Exchange,
			Pair:      payload.Pair,
			Period:    payload.Period,
			Timestamp: current.Unix(),
		}
		c, err := client.cache.Get(key)
		if errors.Is(err, gcache.KeyNotFoundError) { // If not present in cache
			missingTimes = append(missingTimes, current)
			continue
		} else if err != nil { // If there is an error
			return nil, nil, err
		}

		// Check if uncomplete
		cd := c.(candlestick.Candlestick)
		if cd.Uncomplete {
			missingTimes = append(missingTimes, current)
			continue
		}

		// Add to list
		if err := list.Set(current, cd); err != nil {
			return nil, nil, err
		}
	}

	// Change to time ranges and return if none
	return list, timeserie.TimeRangesFromMissingTimes(payload.Period.Duration(), missingTimes), nil
}

func (client *CachedClient) preemptiveAsyncLoading(ctx context.Context, payload ReadCandlesticksPayload) {
	// Round down payload start and end
	payload.Start = payload.End
	payload.End = utils.ToReference(payload.End.Add(payload.Period.Duration() * time.Duration(client.params.PreLoadingSize)))
	payload.Limit = 0

	// Get present and missing times
	_, missingTimeRanges, err := client.presentAndmissingTimes(payload)
	if err != nil {
		telemetry.L(ctx).Errorf("Error while counting missing in preemptive loading: %v", err)
		return
	}

	// Return if there is no need to load
	if len(missingTimeRanges) < client.params.PreLoadingSize/2 {
		return
	}

	telemetry.L(ctx).Debug("Preemptive loading activated")

	// Download missing times
	if _, err = client.downloadMissing(ctx, payload, missingTimeRanges); err != nil {
		telemetry.L(ctx).Errorf("Error while downloading missing in preemptive loading: %v", err)
	}
}

func (client *CachedClient) downloadMissing(ctx context.Context, payload ReadCandlesticksPayload, missingTimeRanges []timeserie.TimeRange) (*candlestick.List, error) {
	// Generate new payload with extended time ranges and limit
	payload.Start = utils.ToReference(missingTimeRanges[0].Start)
	payload.End = utils.ToReference(missingTimeRanges[len(missingTimeRanges)-1].End.Add(payload.Period.Duration() * time.Duration(client.params.PreLoadingSize)))
	if payload.Limit > 0 {
		payload.Limit += uint(client.params.PreLoadingSize)
	}

	// Get missing times
	telemetry.L(ctx).Debugf("Requesting from %s to %s", *payload.Start, *payload.End)
	missing, err := client.controller.Read(ctx, payload)
	if err != nil {
		return nil, err
	}

	// Add missing times to cache
	telemetry.L(ctx).Debugf("Adding %d missing times to cache", missing.Len())
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

	return missing, nil
}

func (client *CachedClient) ServiceInfo(ctx context.Context) (client.ServiceInfo, error) {
	return client.controller.ServiceInfo(ctx)
}

func (client *CachedClient) Close(ctx context.Context) {
	client.cache.Purge()
	client.controller.Close(ctx)
}
