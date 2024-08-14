package cache

import (
	"context"
	"errors"
	"time"

	"cryptellation/internal/adapters/telemetry"
	common "cryptellation/pkg/client"
	"cryptellation/pkg/models/timeserie"
	"cryptellation/pkg/utils"

	client "cryptellation/svc/candlesticks/clients/go"
	"cryptellation/svc/candlesticks/pkg/candlestick"
	"cryptellation/svc/candlesticks/pkg/period"

	"github.com/bluele/gcache"
)

type cacheKey struct {
	Exchange  string
	Pair      string
	Period    period.Symbol
	Timestamp int64
}

type cache struct {
	client   client.Client
	cache    gcache.Cache
	settings struct {
		maxSize                       int
		preLoadingSize                int
		preemptiveAsyncLoadingEnabled bool
	}
}

const (
	DefaultMaxSize              = 100000
	DefaultPreLoadingAfterSize  = 200
	DefaultPreLoadingBeforeSize = 0
)

func New(client client.Client, options ...option) client.Client {
	var c cache

	// Set client and default params
	c.client = client
	c.settings.maxSize = DefaultMaxSize
	c.settings.preLoadingSize = DefaultPreLoadingAfterSize
	c.settings.preemptiveAsyncLoadingEnabled = true

	// Execute options
	for _, option := range options {
		option(&c)
	}

	// Set cache
	c.cache = gcache.New(c.settings.maxSize).LRU().Build()

	return &c
}

func (cache *cache) Read(ctx context.Context, payload client.ReadCandlesticksPayload) (*candlestick.List, error) {
	// Check required fields for caching
	if payload.End == nil {
		payload.End = utils.ToReference(time.Now())
	}
	if payload.Start == nil {
		start := payload.End.Add(-payload.Period.Duration() * time.Duration(cache.settings.preLoadingSize))
		payload.Start = &start
	}

	// Round down payload start and end
	payload.Start = utils.ToReference(payload.Period.RoundTime(*payload.Start))
	payload.End = utils.ToReference(payload.Period.RoundTime(*payload.End))

	// Get present and missing times
	list, tr, err := cache.presentAndmissingTimes(payload)
	if err != nil {
		return nil, err
	}

	if len(tr) == 0 {
		telemetry.L(ctx).Debug("No missing times in cache")
		return list, nil
	}
	telemetry.L(ctx).Debugf("Missing times in cache: %v", tr)

	// Download missing times
	missing, err := cache.downloadMissing(ctx, payload, tr)
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

	if cache.settings.preemptiveAsyncLoadingEnabled && missing.Len() == 0 {
		go cache.preemptiveAsyncLoading(ctx, payload)
	}

	return extract, nil
}

func (cache *cache) presentAndmissingTimes(payload client.ReadCandlesticksPayload) (present *candlestick.List, missing []timeserie.TimeRange, err error) {
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
		c, err := cache.cache.Get(key)
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

func (cache *cache) preemptiveAsyncLoading(ctx context.Context, payload client.ReadCandlesticksPayload) {
	// Round down payload start and end
	payload.Start = payload.End
	payload.End = utils.ToReference(payload.End.Add(payload.Period.Duration() * time.Duration(cache.settings.preLoadingSize)))
	payload.Limit = 0

	// Get present and missing times
	_, missingTimeRanges, err := cache.presentAndmissingTimes(payload)
	if err != nil {
		telemetry.L(ctx).Errorf("Error while counting missing in preemptive loading: %v", err)
		return
	}

	// Return if there is no need to load
	if len(missingTimeRanges) < cache.settings.preLoadingSize/2 {
		return
	}

	telemetry.L(ctx).Debug("Preemptive loading activated")

	// Download missing times
	if _, err = cache.downloadMissing(ctx, payload, missingTimeRanges); err != nil {
		telemetry.L(ctx).Errorf("Error while downloading missing in preemptive loading: %v", err)
	}
}

func (cache *cache) downloadMissing(ctx context.Context, payload client.ReadCandlesticksPayload, missingTimeRanges []timeserie.TimeRange) (*candlestick.List, error) {
	// Generate new payload with extended time ranges and limit
	payload.Start = utils.ToReference(missingTimeRanges[0].Start)
	payload.End = utils.ToReference(missingTimeRanges[len(missingTimeRanges)-1].End.Add(payload.Period.Duration() * time.Duration(cache.settings.preLoadingSize)))
	if payload.Limit > 0 {
		payload.Limit += uint(cache.settings.preLoadingSize)
	}

	// Get missing times
	telemetry.L(ctx).Debugf("Requesting from %s to %s", *payload.Start, *payload.End)
	missing, err := cache.client.Read(ctx, payload)
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
		return false, cache.cache.Set(key, c)
	}); err != nil {
		return nil, err
	}

	return missing, nil
}

func (cache *cache) ServiceInfo(ctx context.Context) (common.ServiceInfo, error) {
	return cache.client.ServiceInfo(ctx)
}

func (cache *cache) Close(ctx context.Context) {
	cache.cache.Purge()
	cache.client.Close(ctx)
}
