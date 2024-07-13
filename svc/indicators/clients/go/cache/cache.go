package cache

import (
	"context"
	"errors"
	"time"

	"github.com/bluele/gcache"
	common "github.com/lerenn/cryptellation/pkg/client"
	"github.com/lerenn/cryptellation/pkg/models/timeserie"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"
	client "github.com/lerenn/cryptellation/svc/indicators/clients/go"
)

type cacheKey struct {
	Exchange     string
	Pair         string
	Period       period.Symbol
	PeriodNumber uint
	PriceType    candlestick.PriceType
	Timestamp    int64
}

type cache struct {
	controller client.Client
	smaCache   gcache.Cache
	settings   struct {
		maxSize              int
		preLoadingAfterSize  int
		preLoadingBeforeSize int
	}
}

const (
	DefaultMaxSize              = 10000
	DefaultPreLoadingAfterSize  = 200
	DefaultPreLoadingBeforeSize = 0
)

func New(controller client.Client, options ...option) client.Client {
	var c cache

	// Set client and default settings
	c.controller = controller
	c.settings.maxSize = DefaultMaxSize
	c.settings.preLoadingAfterSize = DefaultPreLoadingAfterSize
	c.settings.preLoadingBeforeSize = DefaultPreLoadingBeforeSize

	// Execute options
	for _, option := range options {
		option(&c)
	}

	// Set cache
	c.smaCache = gcache.New(c.settings.maxSize).LRU().Build()

	return &c

}

func (client *cache) SMA(ctx context.Context, payload client.SMAPayload) (*timeserie.TimeSerie[float64], error) {
	list := timeserie.New[float64]()

	// Round down payload start and end
	payload.Start, payload.End = payload.Period.RoundInterval(&payload.Start, &payload.End)

	// Get missing times
	missingTimes := make([]time.Time, 0, payload.Period.CountBetweenTimes(payload.Start, payload.End))
	for current := payload.Start; !current.After(payload.End); current = current.Add(payload.Period.Duration()) {
		key := cacheKey{
			Exchange:     payload.Exchange,
			Pair:         payload.Pair,
			Period:       payload.Period,
			PeriodNumber: payload.PeriodNumber,
			PriceType:    payload.PriceType,
			Timestamp:    current.Unix(),
		}
		p, err := client.smaCache.Get(key)
		if errors.Is(err, gcache.KeyNotFoundError) { // If not present in cache
			missingTimes = append(missingTimes, current)
			continue
		} else if err != nil { // If there is an error
			return nil, err
		}

		// Check that it is not now
		if current.Equal(payload.Period.RoundTime(time.Now())) {
			missingTimes = append(missingTimes, current)
			continue
		}

		// Add to list
		_ = list.Set(current, p.(float64))
	}

	// Change to time ranges and return if none
	tr := timeserie.TimeRangesFromMissingTimes(payload.Period.Duration(), missingTimes)
	if len(tr) == 0 {
		return list, nil
	}

	// Generate new payload with extended time ranges
	newPayload := payload
	newPayload.Start = tr[0].Start.Add(-payload.Period.Duration() * time.Duration(client.settings.preLoadingBeforeSize))
	newPayload.End = tr[len(tr)-1].End.Add(payload.Period.Duration() * time.Duration(client.settings.preLoadingAfterSize))

	// Get missing times
	missing, err := client.controller.SMA(ctx, newPayload)
	if err != nil {
		return nil, err
	}

	// Add missing times to cache
	if err := missing.Loop(func(t time.Time, p float64) (bool, error) {
		key := cacheKey{
			Exchange:     payload.Exchange,
			Pair:         payload.Pair,
			Period:       payload.Period,
			PeriodNumber: payload.PeriodNumber,
			PriceType:    payload.PriceType,
			Timestamp:    t.Unix(),
		}
		return false, client.smaCache.Set(key, p)
	}); err != nil {
		return nil, err
	}

	// Merge missing times
	if err := list.Merge(*missing, &timeserie.MergeOptions{}); err != nil {
		return nil, err
	}

	// Exctract only requested
	return list.Extract(payload.Start, payload.End, 0), nil
}

func (client *cache) ServiceInfo(ctx context.Context) (common.ServiceInfo, error) {
	return client.controller.ServiceInfo(ctx)
}

func (client *cache) Close(ctx context.Context) {
	client.smaCache.Purge()
	client.controller.Close(ctx)
}
